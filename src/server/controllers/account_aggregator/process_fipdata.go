package account_aggregator

import (
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/models"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func processAndSaveFipDataCollection(allFipData []fipDataCollection,
	userConsent models.UserConsents) error {

	name := getHolderField(allFipData, "Name")
	dob, _ := time.Parse("02-01-2006", getHolderField(allFipData, "Dob"))
	panCard := getHolderField(allFipData, "Pan")
	ckycCompliance, _ := strconv.ParseBool(getHolderField(allFipData, "CkycCompliance"))

	planInformation := models.UserScores{
		UserConsentId:   userConsent.Id,
		Name:            name,
		DateOfBirth:     dob,
		PanCard:         panCard,
		CkycCompliance:  ckycCompliance,
		AgeScore:        getAgeScore(dob),
		WealthScore:     getWealthPlanScore(allFipData),
		DebtScore:       getDebtPlanScore(allFipData),
		MedicalScore:    getMedicalPlanScore(allFipData),
		InvestmentScore: getInvestmentPlanScore(allFipData),
		PensionScore:    getPensionPlanScore(allFipData),
		FamilyScore:     getFamilyPlanScore(allFipData),
		ChildrenScore:   getChildrenPlanScore(allFipData),
		MotorScore:      getMotorPlanScore(allFipData),
		TermScore:       getTermPlanScore(allFipData),
		TravelScore:     getTravelPlanScore(allFipData),
		AllScore:        getAllPlanScore(allFipData),
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
	return nil
}

func getHolderField(allFipData []fipDataCollection, field string) string {
	for _, fipData := range allFipData {
		for _, fipDataItem := range fipData.FipData {
			// Problem in the Setu API's new version
			for _, holder := range fipDataItem.Account.Profile.Holders.Holder {
				relfectedHolder := reflect.ValueOf(holder)
				fieldValue := reflect.Indirect(relfectedHolder).FieldByName(field)
				if fieldValue.String() != "" {
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
			if fipDataItem.Account.Type == "insurance_policies" {
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
