package insurance

import (
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/controllers"
	"insurance-ng/src/server/models"
	"regexp"

	"github.com/google/uuid"
)

func getUserDataRecord(userID string) (responseData getUserDataResponse, err error) {
	userConsentCh := make(chan userConsentChResp, 1)
	go getUserConsent(userID, userConsentCh)
	insuranceAvailableCh := make(chan insuranceChResp, 1)
	go getAllInsuranceOffers(insuranceAvailableCh)

	insuranceAvailable := <-insuranceAvailableCh
	if insuranceAvailable.err != nil {
		err = insuranceAvailable.err
		return
	}

	userConsent := <-userConsentCh
	if userConsent.err != nil {
		err = userConsent.err
		return
	}

	if userConsent.result.ArtefactStatus != models.ArtefactStatusReady &&
		userConsent.result.ArtefactStatus != models.ArtefactStatusActive {
		return getUserDataResponse{Status: true}, nil
	}

	if !userConsent.result.DataFetched {
		return getUserDataResponse{Status: false}, nil
	}

	userPlanScoresCh := make(chan userScoreChResp, 1)
	userExistingInsuranceCh := make(chan userInsurancesChResp, 1)
	go getUserPlanScore(userConsent.result.ID, userPlanScoresCh)
	go getUserInsurances(userConsent.result.CustomerID, userExistingInsuranceCh)

	userScore := <-userPlanScoresCh
	if userScore.err != nil {
		err = insuranceAvailable.err
		return
	}

	userExistingInsurance := <-userExistingInsuranceCh
	if userExistingInsurance.err != nil {
		err = userConsent.err
		return
	}

	var insuranceOffered []insuranceOffers

	for _, insuranceInfo := range insuranceAvailable.result {
		existingInsurance := models.UserInsurance{}
		for _, userInsurances := range userExistingInsurance.result {
			if userInsurances.Type == insuranceInfo.Type {
				existingInsurance = *userInsurances
				break
			}
		}

		insuranceScore := getScoreForType(userScore.result, insuranceInfo.Type)
		insuranceOffered = append(insuranceOffered, insuranceOffers{
			ID:                insuranceInfo.ID,
			Title:             insuranceInfo.Title,
			Description:       insuranceInfo.Description,
			AccountID:         existingInsurance.AccountID,
			Score:             insuranceScore,
			CurrentPremium:    existingInsurance.Premium,
			CurrentCover:      existingInsurance.Cover,
			Clauses:           insuranceInfo.Clauses,
			CurrentClauses:    existingInsurance.Clauses,
			OfferedPremium:    getOfferPremium(insuranceInfo.Premium, float64(insuranceScore)),
			OfferedCover:      getOfferCover(insuranceInfo.Cover, float64(insuranceScore)),
			YoyDeductionRate:  insuranceInfo.YoyDeductionRate,
			IsInsuranceNgAcct: existingInsurance.IsInsuranceNgAcct,
			IsActive:          existingInsurance.IsActive,
			IsClaimed:         existingInsurance.IsClaimed,
			Type:              insuranceInfo.Type,
		})
	}

	regexPattern := regexp.MustCompile("@.*$")
	responseData = getUserDataResponse{
		Status: true,
		Data: userData{
			Name:              userScore.result.Name,
			DateOfBirth:       userScore.result.DateOfBirth,
			Pancard:           userScore.result.Pancard,
			Phone:             regexPattern.ReplaceAllString(userConsent.result.CustomerID, ""),
			SharedDataSources: userScore.result.SharedDataSources,
			CkycCompliance:    userScore.result.CkycCompliance,
			AgeScore:          userScore.result.AgeScore,
			WealthScore:       userScore.result.WealthScore,
			DebtScore:         userScore.result.DebtScore,
			InvestmentScore:   userScore.result.InvestmentScore,
			InsuranceOffers:   insuranceOffered,
		},
		Error: nil,
	}

	return
}

func getUserInsurances(pancard string, userExistingInsuranceCh chan userInsurancesChResp) {
	var userExistingInsurance []*models.UserInsurance
	result := config.Database.Where("customer_id = ?", pancard).Find(&userExistingInsurance)

	userExistingInsuranceCh <- userInsurancesChResp{
		result: userExistingInsurance,
		err:    result.Error,
	}
}

func getUserPlanScore(consentID uuid.UUID, userPlanScoreCh chan userScoreChResp) {
	var userPlanScore *models.UserScores
	result := config.Database.Where("user_consent_id = ?", consentID).Take(&userPlanScore)

	userPlanScoreCh <- userScoreChResp{
		result: userPlanScore,
		err:    result.Error,
	}
}

func getUserConsent(userID string, userConsentChannel chan userConsentChResp) {
	var userConsent *models.UserConsents
	result := config.Database.Where("user_id = ?", userID).Take(&userConsent)

	userConsentChannel <- userConsentChResp{
		result: userConsent,
		err:    result.Error,
	}
}

func getAllInsuranceOffers(insuranceAvailableChannel chan insuranceChResp) {
	var insurance []*models.Insurance
	result := config.Database.Find(&insurance)

	insuranceAvailableChannel <- insuranceChResp{
		result: insurance,
		err:    result.Error,
	}
}

func createInsuranceRecord(userID string, insuranceID uuid.UUID) (err error) {
	var userConsent models.UserConsents
	if result := config.Database.Model(&models.UserConsents{}).Where("user_id = ?",
		userID).Take(&userConsent); result.Error != nil {
		err = result.Error
		return
	}

	var userScore models.UserScores
	if result := config.Database.Model(&models.UserScores{}).Where("user_consent_id = ?",
		userConsent.ID).Take(&userScore); result.Error != nil {
		err = result.Error
		return
	}

	var insurance models.Insurance
	if result := config.Database.Model(&models.Insurance{}).Where("id = ?",
		insuranceID).Take(&insurance); result.Error != nil {
		err = result.Error
		return
	}

	insuranceScore := float64(getScoreForType(&userScore, insurance.Type))
	insuranceAcctID := ""
	insuranceActivate := insuranceScore > PreApprovedBar
	if insuranceActivate {
		insuranceAcctID = controllers.GetRandomString(10)
	}
	newInsurance := models.UserInsurance{
		Type:              insurance.Type,
		Premium:           getOfferPremium(insurance.Premium, insuranceScore),
		Cover:             getOfferCover(insurance.Cover, insuranceScore),
		IsActive:          insuranceActivate,
		AccountID:         insuranceAcctID,
		CustomerID:        userConsent.CustomerID,
		Clauses:           insurance.Clauses,
		IsClaimed:         false,
		IsInsuranceNgAcct: true,
	}

	if config.Database.Where("customer_id = ?", userConsent.CustomerID).Where("type = ?",
		insurance.Type).Updates(&newInsurance).RowsAffected == 0 {
		if result := config.Database.Create(&newInsurance); result.Error != nil {
			err = result.Error
			return
		}
	}

	return
}

func initiateInsuranceClaim(userID string, insuranceID uuid.UUID) (err error) {
	var userConsent models.UserConsents
	if result := config.Database.Model(&models.UserConsents{}).Where("user_id = ?",
		userID).Take(&userConsent); result.Error != nil {
		err = result.Error
		return
	}

	var insurance models.Insurance
	if result := config.Database.Model(&models.Insurance{}).Where("id = ?",
		insuranceID).Take(&insurance); result.Error != nil {
		err = result.Error
		return
	}

	result := config.Database.Model(&models.UserInsurance{}).Where("customer_id = ?",
		userConsent.CustomerID).Where("type = ?", insurance.Type).Update("is_claimed", true)
	if result.Error != nil {
		err = result.Error
		return
	}

	return
}
