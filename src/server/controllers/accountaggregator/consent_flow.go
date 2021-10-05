package accountaggregator

import (
	"encoding/json"
	"fmt"
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/models"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/google/uuid"
)

// Create Consent
func sendCreateConsentReqToAcctAggregator(phone string) (consentResponse setuCreateConsentResponse,
	consentExpire time.Time, err error) {
	customerID := fmt.Sprintf("%s@setu-aa", phone)
	consentUUID := uuid.New()
	startTime := time.Now()
	endTime := time.Now().Add(time.Minute * 15)

	consentBody := createConsentBody(consentUUID, startTime, endTime, customerID)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodRS256, consentBody)
	setuRequestBody, err := json.Marshal(consentBody)
	if err != nil {
		return
	}

	respBytes, err := sendRequestToSetu(SetuAPICreateConsentPath, "POST", setuRequestBody, jwtToken)
	if err != nil {
		return
	}

	if err = json.Unmarshal(respBytes, &consentResponse); err != nil {
		return
	}

	consentExpire = endTime
	return
}

func createConsentBody(uuid uuid.UUID, startTime time.Time, endTime time.Time, customerID string) setuCreateConsentRequest {

	// Time Hack:
	// Currently, Setu API is not following the RFC3339 Correctly,
	// Hence for the time being, we manually converting dates.
	// startTime := time.Now().Format(time.RFC3339)
	// endTime := time.Now().Add(time.Minute * 15).Format(time.RFC3339)
	startTime3339 := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.153Z",
		startTime.Year(), startTime.Month(), startTime.Day(),
		startTime.Hour(), startTime.Minute(), startTime.Second())
	fromTime3339 := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.153Z",
		endTime.Year()-5, endTime.Month(), endTime.Day(),
		endTime.Hour(), endTime.Minute(), endTime.Second())
	endTime3339 := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.153Z",
		endTime.Year(), endTime.Month(), endTime.Day(),
		endTime.Hour(), endTime.Minute(), endTime.Second())

	const dataEndpointsExplaination = `To provide the lowest prices, we need your financial information to understand you. We take the data points like personal information (including name, date of birth, address, pancard & existing insurance plans), deposit account transactions, current balance summary, insurance accounts and transactions, investment plans and debt to understand your lifestyle and plans that would best suit you. You have the right to revoke/request for data deletion at any point in the future. Sharing bank account and insurance account is a minimum requirement.`

	requestBody := setuCreateConsentRequest{
		Ver:       "1.0",
		Timestamp: startTime3339,
		Txnid:     uuid,
		ConsentDetail: consentDetails{
			ConsentStart:  startTime3339,
			ConsentExpiry: endTime3339,
			ConsentMode:   ConsentModeView,
			FetchType:     FetchTypeOnetime,
			ConsentTypes:  []string{ConsentTypesProfile, ConsentTypesSummary, ConsentTypesTransaction},
			FiTypes: []string{
				FiTypesDeposit, FiTypesInsurancePolicies, FiTypesCreditCard,
				FiTypesTermDeposit, FiTypesRecurringDeposit, FiTypesPPF, FiTypesNPS,
			},
			DataConsumer: idType{
				ID: "FIU",
			},
			Customer: idType{
				ID: customerID,
			},
			Purpose: purpose{
				Code:   fmt.Sprint(PurposeOneTime),
				RefURI: "https://api.rebit.org.in/aa/purpose/105.xml",
				Text:   dataEndpointsExplaination,
				Category: purposeCategory{
					Type: "string",
				},
			},
			FIDataRange: fIDataRange{From: fromTime3339, To: endTime3339},
			DataLife:    dataTimeRange{Unit: "DAY", Value: 1},
			Frequency:   dataTimeRange{Unit: "DAY", Value: 1},
			DataFilter:  []dataFilter{},
		},
	}

	return requestBody
}

func addOrUpdateConsentToDb(userID string, consent setuCreateConsentResponse,
	expiry time.Time) (userConsent models.UserConsents, err error) {

	userConsent = models.UserConsents{
		UserID:         userID,
		CustomerID:     consent.Customer.ID,
		Expire:         expiry,
		ConsentStatus:  models.EmptyColumn,
		ArtefactStatus: models.ArtefactStatusPending,
		DataFetched:    false,
		SignedConsent:  models.EmptyColumn,
		ConsentHandle:  consent.ConsentHandle,
		ConsentID:      uuid.Nil,
	}

	if err = deleteUserConsent(userConsent); err != nil {
		return models.UserConsents{}, err
	}

	if result := config.Database.Create(&userConsent); result.Error != nil {
		return models.UserConsents{}, result.Error
	}
	return
}

func updateUserConsentForStatusChange(consent consentNotifierStatus) error {
	var userConsent models.UserConsents
	if consent.ConsentStatus == models.ConsentStatusRevoked {
		userConsent = models.UserConsents{
			SignedConsent:  models.EmptyColumn,
			ConsentID:      consent.ConsentID,
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

	consentID, err := getUserArtefactStatus(userConsent.UserID, userConsent.ConsentHandle)
	if err != nil {
		return
	}
	if err = fetchSignedConsent(userConsent.UserID, consentID); err != nil {
		return
	}
	if err = createAndSaveSessionDetails(userConsent.UserID); err != nil {
		return
	}

	return
}

func getUserArtefactStatus(userID string, consentHandle uuid.UUID) (consentID uuid.UUID, err error) {
	urlPath := fmt.Sprintf(SetuAPIConsentStatusPath, consentHandle)
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

	consentID = consentStatus.ConsentStatus.ID
	return
}

func fetchSignedConsent(userID string, consentID uuid.UUID) (err error) {
	urlPath := fmt.Sprintf(SetuAPIConsentSignedPath, consentID)
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
		ConsentID:      signedConsent.ConsentID,
	}

	updatedRow := config.Database.Model(&models.UserConsents{}).Where("user_id = ?",
		userID).Updates(&updatedConsent)
	return updatedRow.Error
}
