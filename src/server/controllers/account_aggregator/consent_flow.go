package account_aggregator

import (
	"encoding/json"
	"fmt"
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/controllers"
	"insurance-ng/src/server/models"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Create Consent
func sendCreateConsentReqToAcctAggregator(phone string) (consentResponse setuCreateConsentResponse,
	consentExpire time.Time, err error) {
	customerId := fmt.Sprintf("%s@setu-aa", phone)
	consentUuid := uuid.New()
	startTime := time.Now()
	endTime := time.Now().Add(time.Minute * 15)

	consentBody := createConsentBody(consentUuid, startTime, endTime, customerId)
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

func createConsentBody(uuid uuid.UUID, startTime time.Time, endTime time.Time, customerId string) setuCreateConsentRequest {

	// TODO: Unfortunatly, there are some bugs in the Setu API as of today, which
	// I have reported, but until those bugs are fixed, we have to comment this and
	// use the hack below.
	// startTime := time.Now().Format(time.RFC3339)
	// endTime := time.Now().Add(time.Minute * 15).Format(time.RFC3339)
	startTimeHack := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.153Z",
		startTime.Year(), startTime.Month(), startTime.Day(),
		startTime.Hour(), startTime.Minute(), startTime.Second())
	fromTimeHack := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.153Z",
		endTime.Year()-5, endTime.Month(), endTime.Day(),
		endTime.Hour(), endTime.Minute(), endTime.Second())
	endTimeHack := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.153Z",
		endTime.Year(), endTime.Month(), endTime.Day(),
		endTime.Hour(), endTime.Minute(), endTime.Second())

	requestBody := setuCreateConsentRequest{
		Ver:       "1.0",
		Timestamp: startTimeHack,
		Txnid:     uuid,
		ConsentDetail: consentDetails{
			ConsentStart:  startTimeHack,
			ConsentExpiry: endTimeHack,
			ConsentMode:   ConsentModeView,
			FetchType:     FetchTypeOnetime,
			ConsentTypes:  []string{ConsentTypesProfile, ConsentTypesSummary, ConsentTypesTransaction},
			FiTypes: []string{
				FiTypesDeposit, FiTypesInsurancePolicies, FiTypesMutualFunds,
				FiTypesTermDeposit, FiTypesRecurringDeposit, FiTypesPPF,
				// FiTypesNPS, FiTypesSIP, FiTypesGovernmentSecrities, FiTypesEquities,
			},
			DataConsumer: idType{
				Id: "FIU",
			},
			Customer: idType{
				Id: customerId,
			},
			Purpose: purpose{
				Code:   fmt.Sprint(PurposeOneTime),
				RefUri: "https://api.rebit.org.in/aa/purpose/105.xml",
				Text:   controllers.DataEndpointsExplainationInParagraph,
				Category: purposeCategory{
					Type: "string",
				},
			},
			FIDataRange: fIDataRange{From: fromTimeHack, To: endTimeHack},
			DataLife:    dataTimeRange{Unit: "DAY", Value: 1},
			Frequency:   dataTimeRange{Unit: "DAY", Value: 1},
			DataFilter:  []dataFilter{},
		},
	}

	return requestBody
}

func addOrUpdateConsentToDb(userId string, consent setuCreateConsentResponse, expiry time.Time) *gorm.DB {
	userConsent := models.UserConsents{
		UserId:         userId,
		CustomerId:     consent.Customer.Id,
		Expire:         expiry,
		ConsentStatus:  models.EmptyColumn,
		ArtefactStatus: models.ArtefactStatusPending,
		DataFetched:    false,
		SignedConsent:  models.EmptyColumn,
		ConsentHandle:  consent.ConsentHandle,
		ConsentId:      uuid.Nil,
	}

	deleteUserConsent(userConsent)
	return config.Database.Create(&userConsent)
}

func updateUserConsentForStatusChange(consent consentNotifierStatus) error {
	var userConsent models.UserConsents
	if consent.ConsentStatus == models.ConsentStatusRevoked {
		userConsent = models.UserConsents{
			SignedConsent:  models.EmptyColumn,
			ConsentId:      consent.ConsentId,
			DataFetched:    false,
			ConsentStatus:  consent.ConsentStatus,
			ArtefactStatus: models.ArtefactStatusDenied,
		}
	} else {
		userConsent = models.UserConsents{
			ConsentStatus: consent.ConsentStatus,
		}
	}
	if result := config.Database.Model(&models.UserConsents{}).Where("consent_handle = ?",
		consent.ConsentHandle).Updates(&userConsent).Take(&userConsent); result.Error != nil {
		return result.Error
	}

	if consent.ConsentStatus == models.ConsentStatusActive {
		return updateSignedConsent(userConsent)
	}
	return nil
}

func updateSignedConsent(userConsent models.UserConsents) (err error) {
	// If the data is already fetched, we are all good.
	// Let's inform the user that data request can be initiated.
	if userConsent.DataFetched || userConsent.SignedConsent != models.EmptyColumn {
		return
	}

	consentId, err := getUserArtefactStatus(userConsent.UserId, userConsent.ConsentHandle)
	if err != nil {
		return
	}
	if err = fetchSignedConsent(userConsent.UserId, consentId); err != nil {
		return
	}
	if err = createAndSaveSessionDetails(userConsent.UserId); err != nil {
		return
	}

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
	fmt.Println(string(respBytes))

	var consentStatus setuConsentStatusResponse
	if err = json.Unmarshal(respBytes, &consentStatus); err != nil {
		return
	}

	consentId = consentStatus.ConsentStatus.Id
	return
}

func fetchSignedConsent(userId string, consentId uuid.UUID) (err error) {
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
		ArtefactStatus: signedConsent.Status,
		SignedConsent:  signedConsent.SignedConsent,
		ConsentId:      signedConsent.ConsentId,
	}

	updatedRow := config.Database.Model(&models.UserConsents{}).Where("user_id = ?", userId).Updates(&updatedConsent)
	return updatedRow.Error
}
