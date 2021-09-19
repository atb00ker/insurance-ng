package account_aggregator

import (
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/models"
	"reflect"
	"strconv"
	"time"
)

func processAndSaveFipDataCollection(allFipData []fipDataCollection,
	userConsent models.UserConsents) (err error) {

	name := getHolderField(allFipData, "Name")
	dob, _ := time.Parse("2006-01-02", getHolderField(allFipData, "Dob"))
	panCard := getHolderField(allFipData, "Pan")
	ckycCompliance, _ := strconv.ParseBool(getHolderField(allFipData, "CkycCompliance"))

	planInformation := models.UserPlanScores{
		UserConsentId:     userConsent.Id,
		Name:              name,
		DateOfBirth:       dob,
		PanCard:           panCard,
		CkycCompliance:    ckycCompliance,
		AgeScore:          getAgeScore(dob),
		MedicalScore:      getMedicalPlanScore(allFipData),
		PensionScore:      getPensionPlanScore(allFipData),
		FamilyScore:       getFamilyPlanScore(allFipData),
		ChildrenScore:     getChildrenPlanScore(allFipData),
		MotorScore:        getMotorPlanScore(allFipData),
		TermScore:         getTermPlanScore(allFipData),
		TravelScore:       getTravelPlanScore(allFipData),
		AllScore:          getAllPlanScore(allFipData),
		MedicalAccountId:  getAccountIdField(allFipData, models.OfferePlansMedicalPlan),
		PensionAccountId:  getAccountIdField(allFipData, models.OfferePlansPensionPlan),
		FamilyAccountId:   getAccountIdField(allFipData, models.OfferePlansFamilyPlan),
		ChildrenAccountId: getAccountIdField(allFipData, models.OfferePlansChildrenPlan),
		MotorAccountId:    getAccountIdField(allFipData, models.OfferePlansMotorPlan),
		TermAccountId:     getAccountIdField(allFipData, models.OfferePlansTermPlan),
		TravelAccountId:   getAccountIdField(allFipData, models.OfferePlansTravelPlan),
		AllAccountId:      getAccountIdField(allFipData, models.OfferePlansAllPlan),
	}

	config.Database.Create(&planInformation)
	config.Database.Model(&models.UserConsents{}).Where("id = ?", userConsent.Id).Update("data_fetched", true)
	return
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
