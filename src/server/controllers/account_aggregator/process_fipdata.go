package account_aggregator

import (
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/controllers/insurance"
	"insurance-ng/src/server/models"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func processAndSaveFipDataCollection(allFipData []fipDataCollection,
	userConsent models.UserConsents) error {

	sharedDataSources := getSharedDataSources(allFipData)
	name := getDepositHolderField(allFipData, "Name")
	dob, _ := time.Parse("02-01-2006", getDepositHolderField(allFipData, "Dob"))
	panCard := getDepositHolderField(allFipData, "Pan")
	ckycCompliance, _ := strconv.ParseBool(getDepositHolderField(allFipData, "CkycCompliance"))

	planInformation := models.UserScores{
		UserConsentId:     userConsent.Id,
		Name:              name,
		DateOfBirth:       dob,
		PanCard:           panCard,
		CkycCompliance:    ckycCompliance,
		SharedDataSources: sharedDataSources,
		AgeScore:          getAgeScore(dob),
		WealthScore:       getWealthPlanScore(allFipData, sharedDataSources),
		DebtScore:         getDebtPlanScore(allFipData, sharedDataSources),
		MedicalScore:      getMedicalPlanScore(allFipData, sharedDataSources),
		InvestmentScore:   getInvestmentScore(allFipData, sharedDataSources),
		PensionScore:      getPensionPlanScore(allFipData, sharedDataSources),
		FamilyScore:       getFamilyPlanScore(allFipData, sharedDataSources),
		ChildrenScore:     getChildrenPlanScore(allFipData, sharedDataSources),
		MotorScore:        getMotorPlanScore(allFipData, sharedDataSources),
		TermScore:         getTermPlanScore(allFipData, sharedDataSources),
		TravelScore:       getTravelPlanScore(allFipData, sharedDataSources),
		AllScore:          getAllPlanScore(allFipData, sharedDataSources),
	}
	// We can to this in go routines to make this section faster.
	if err := saveExistingInsuranceInformation(allFipData, userConsent.Id); err != nil {
		return err
	}
	if result := config.Database.Create(&planInformation); result.Error != nil {
		return result.Error
	}
	if result := config.Database.Model(&models.UserConsents{}).Where("id = ?",
		userConsent.Id).Update("data_fetched", true); result.Error != nil {
		return result.Error
	}

	insurance.WaitForProcessing <- userConsent.UserId
	return nil
}

func getSharedDataSources(allFipData []fipDataCollection) int16 {
	var records int16 = 0
	basicRequirements := map[string]bool{
		strings.ToLower(FiTypesInsurancePolicies): false,
		strings.ToLower(FiTypesDeposit):           false,
	}

	for _, fipData := range allFipData {
		for _, fipDataItem := range fipData.FipData {
			delete(basicRequirements, fipDataItem.Account.Type)
			records++
		}
	}

	if len(basicRequirements) != 0 {
		// Basic account are not shared, we will not consider it a
		// valid consent request.
		return 0
	}

	return records
}

func getDepositHolderField(allFipData []fipDataCollection, field string) string {
	for _, fipData := range allFipData {
		for _, fipDataItem := range fipData.FipData {
			if fipDataItem.Account.Type == strings.ToLower(FiTypesDeposit) {
				// Problem in the Setu API's new version
				for _, holder := range fipDataItem.Account.Profile.Holders.Holder {
					relfectedHolder := reflect.ValueOf(holder)
					fieldValue := reflect.Indirect(relfectedHolder).FieldByName(field)
					// Return as soon as we find a name
					// We can match for mobile number to ensure we can
					// picking the correct person's defails but for
					// mock data, that is not required.
					return fieldValue.String()
				}
			}
		}
	}
	return ""
}

func saveExistingInsuranceInformation(allFipData []fipDataCollection, consendId uuid.UUID) error {
	for _, fipData := range allFipData {
		for _, fipDataItem := range fipData.FipData {
			if fipDataItem.Account.Type == strings.ToLower(FiTypesInsurancePolicies) {
				premium, _ := strconv.ParseFloat(fipDataItem.Account.Summary.PremiumAmount, 32)
				cover, _ := strconv.ParseFloat(fipDataItem.Account.Summary.CoverAmount, 32)
				insurance := models.UserInsurance{
					UserConsentId: consendId,
					Type:          fipDataItem.Account.Summary.PolicyType,
					Premium:       premium,
					Cover:         cover,
					Clauses:       models.InsuranceApnaMockedClauses,
					AccountId:     fipDataItem.Account.MaskedAccNumber,
				}
				result := config.Database.Create(&insurance)
				if result.Error != nil {
					return result.Error
				}
			}
		}
	}
	return nil
}
