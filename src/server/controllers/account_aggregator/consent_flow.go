package account_aggregator

import (
	"encoding/json"
	"fmt"
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/models"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func sendCreateConsentReqToAcctAggregator(phone string) (consentResponse setuCreateConsentResponse,
	consentExpire time.Time, err error) {
	// TODO: Unfortunatly, there are some bugs in the Setu API as of today, which
	// I have reported, but until those bugs are fixed, we have to comment this and
	// use the hack below.
	// startTime := time.Now().Format(time.RFC3339)
	// endTime := time.Now().Add(time.Minute * 15).Format(time.RFC3339)
	startTime := time.Now()
	endTime := time.Now().Add(time.Minute * 15)
	startTimeHack := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.153Z",
		startTime.Year(), startTime.Month(), startTime.Day(),
		startTime.Hour(), startTime.Minute(), startTime.Second())
	endTimeHack := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.153Z",
		endTime.Year(), endTime.Month(), endTime.Day(),
		endTime.Hour(), endTime.Minute(), endTime.Second())
	customerId := fmt.Sprintf("%s@setu-aa", phone)
	consentUuid := uuid.New()
	consentBody := createConsentBody(consentUuid, startTimeHack, endTimeHack, customerId)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, consentBody)
	setuRequestBody, err := json.Marshal(consentBody)
	if err != nil {
		return
	}

	respBytes, err := sendRequestToSetu(SetuApiCreateConsentPath, "POST", setuRequestBody, jwtToken)
	if err != nil {
		return
	}

	if err = json.Unmarshal(respBytes, &consentResponse); err != nil {
		return
	}

	consentExpire = endTime
	return
}

func createConsentBody(uuid uuid.UUID, startTime string, endTime string,
	customerId string) setuCreateConsentRequest {

	requestBody := setuCreateConsentRequest{
		Ver:       "1.0",
		Timestamp: startTime,
		Txnid:     uuid,
		ConsentDetail: consentDetails{
			ConsentStart:  startTime,
			ConsentExpiry: endTime,
			ConsentMode:   ConsentModeView,
			FetchType:     FetchTypeOnetime,
			ConsentTypes:  []string{ConsentTypesProfile, ConsentTypesSummary, ConsentTypesTransaction},
			FiTypes:       []string{FiTypesDeposit, FiTypesInsurancePolicies, FiTypesMutualFunds},
			DataConsumer: idType{
				Id: "FIU",
			},
			Customer: idType{
				Id: customerId,
			},
			Purpose: purpose{
				Code:   fmt.Sprint(PurposeOneTime),
				RefUri: "https://api.rebit.org.in/aa/purpose/105.xml",
				Text:   "Management System",
				Category: purposeCategory{
					Type: "string",
				},
			},
			FIDataRange: fIDataRange{From: "1947-08-15T00:00:00.153Z", To: endTime},
			DataLife:    dataTimeRange{Unit: "DAY", Value: 1},
			Frequency:   dataTimeRange{Unit: "DAY", Value: 1},
			DataFilter:  []dataFilter{},
		},
	}

	return requestBody
}

func addOrUpdateConsentToDb(userId string, consent setuCreateConsentResponse, expiry time.Time) *gorm.DB {
	userConsent := models.UserConsents{
		UserId:        userId,
		CustomerId:    consent.Customer.Id,
		Expire:        expiry,
		Status:        models.ConsentPending,
		SignedConsent: "-",
		UserData:      "-",
		ConsentHandle: consent.ConsentHandle,
		ConsentId:     uuid.Nil,
	}
	updatedRow := config.Database.Model(&models.UserConsents{}).Where("user_id = ?", userId).Updates(&userConsent)
	if updatedRow.RowsAffected == 0 {
		return config.Database.Create(&userConsent)
	}

	return updatedRow
}

func getUserConsent(userId string) (userConsent models.UserConsents, err error) {
	result := config.Database.Model(&models.UserConsents{}).Where("user_id = ?", userId).Take(&userConsent)
	err = result.Error
	return
}

func getUserArtefactStatus(userId string, consentHandle uuid.UUID) (consentId uuid.UUID, err error) {
	urlPath := fmt.Sprintf(SetuApiConsentStatusPath, consentHandle)
	jwtPayload := setuConsentStatusRequest{Path: urlPath}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtPayload)
	respBytes, err := sendRequestToSetu(urlPath, "GET", []byte{}, jwtToken)
	if err != nil {
		return
	}

	var consentStatus setuConsentStatusResponse
	if err = json.Unmarshal(respBytes, &consentStatus); err != nil {
		return
	}

	consentId = consentStatus.ConsentStatus.Id
	return
}

func fetchSignedConsent(userId string, consentId uuid.UUID) (status string, err error) {
	urlPath := fmt.Sprintf(SetuApiConsentSignedPath, consentId)
	jwtPayload := setuConsentStatusRequest{Path: urlPath}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtPayload)
	respBytes, err := sendRequestToSetu(urlPath, "GET", []byte{}, jwtToken)
	if err != nil {
		return
	}

	var signedConsent setuSignedConsentResponse
	if err = json.Unmarshal(respBytes, &signedConsent); err != nil {
		return
	}

	updatedConsent := models.UserConsents{
		Status:        signedConsent.Status,
		SignedConsent: signedConsent.SignedConsent,
		ConsentId:     signedConsent.ConsentId,
	}

	updatedRow := config.Database.Model(&models.UserConsents{}).Where("user_id = ?", userId).Updates(&updatedConsent)
	return updatedConsent.Status, updatedRow.Error
}
