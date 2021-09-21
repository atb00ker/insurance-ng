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
	dob, _ := time.Parse("2006-01-02", getHolderField(allFipData, "Dob"))
	panCard := getHolderField(allFipData, "Pan")
	ckycCompliance, _ := strconv.ParseBool(getHolderField(allFipData, "CkycCompliance"))

	planInformation := models.UserPlanScores{
		UserConsentId:  userConsent.Id,
		Name:           name,
		DateOfBirth:    dob,
		PanCard:        panCard,
		CkycCompliance: ckycCompliance,
		AgeScore:       getAgeScore(dob),
		MedicalScore:   getMedicalPlanScore(allFipData),
		PensionScore:   getPensionPlanScore(allFipData),
		FamilyScore:    getFamilyPlanScore(allFipData),
		ChildrenScore:  getChildrenPlanScore(allFipData),
		MotorScore:     getMotorPlanScore(allFipData),
		TermScore:      getTermPlanScore(allFipData),
		TravelScore:    getTravelPlanScore(allFipData),
		AllScore:       getAllPlanScore(allFipData),
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

func getAccountIdField(allFipData []fipDataCollection, planName string) string {
	for _, fipData := range allFipData {
		for _, fipDataItem := range fipData.FipData {
			if fipDataItem.Account.Type == "insurance" &&
				fipDataItem.Account.Summary.PolicyType == planName {
				return fipDataItem.Account.MaskedAccNumber
			}
		}
	}
	return ""
}

func getHolderField(allFipData []fipDataCollection, field string) string {
	for _, fipData := range allFipData {
		for _, fipDataItem := range fipData.FipData {
			for _, holder := range fipDataItem.Account.Profile.Holders {
				relfectedHolder := reflect.ValueOf(holder.Holder)
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
			if fipDataItem.Account.Type == "insurance" {
				insurance := models.UserExistingInsurance{
					UserConsentId: consendId,
					Type:          fipDataItem.Account.Summary.PolicyType,
					Premium:       fipDataItem.Account.Summary.PremiumAmount,
					Cover:         fipDataItem.Account.Summary.CoverAmount,
					AccountId:     getAccountIdField(allFipData, fipDataItem.Account.Summary.PolicyType),
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
