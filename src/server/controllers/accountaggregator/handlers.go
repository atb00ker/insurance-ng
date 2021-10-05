package accountaggregator

import (
	"encoding/json"
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/controllers"
	"insurance-ng/src/server/models"
	"net/http"

	"github.com/google/uuid"
)

// URLs for account aggregator endpoints
const (
	URLCreateConsent           string = "/api/v1/account_aggregator/consent/"
	URLConsentNotificationMock string = "/api/v1/account_aggregator/Mock/Consent/Notification/"
	URLConsentNotification     string = "/api/v1/account_aggregator/Consent/Notification"
	URLArtefactNotification    string = "/api/v1/account_aggregator/FI/Notification"
)

// CreateConsentRequest creates a new AA Consent request
func CreateConsentRequest(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userID, err := controllers.GetUserIdentifier(response, request)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	decoder := json.NewDecoder(request.Body)
	var requestJSON createConsentRequestInput
	if err := decoder.Decode(&requestJSON); err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	consentResponse, expiry, err := sendCreateConsentReqToAcctAggregator(requestJSON.Phone)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	if _, err = addOrUpdateConsentToDb(userID, consentResponse, expiry); err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	respMessage, _ := json.Marshal(createConsentResponseOutput{
		ConsentHandle: consentResponse.ConsentHandle.String(),
	})

	response.Write(respMessage)
}

// ConsentNotification initiates state update on consent notification
func ConsentNotification(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(request.Body)
	var requestJSON setuConsentNotificationRequest
	if err := decoder.Decode(&requestJSON); err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	if err := updateUserConsentForStatusChange(requestJSON.ConsentStatusNotification); err != nil {
		handleNotificationError(response, err)
		return
	}

	clientAPI, requestJWS, setuResponseBody, err := sendResponseToSetuNotification()
	if err != nil {
		handleNotificationError(response, err)
		return
	}

	response.Header().Set("client_api_key", clientAPI)
	response.Header().Set("x-jws-signature", requestJWS)
	response.Write(setuResponseBody)
}

// ConsentNotificationMock returns mocked notification response
func ConsentNotificationMock(response http.ResponseWriter, request *http.Request) {
	// Notification Hack:
	// This is just a mock endpoint.
	// Currently the Setu notifications are not triggered
	// sometimes, which can be a big problem in the
	// hackathon, hence, for the time being, we start the
	// notification steps here. Not required in production.
	userID, err := controllers.GetUserIdentifier(response, request)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	var userConsent models.UserConsents
	if result := config.Database.Model(&models.UserConsents{}).Where("user_id = ?",
		userID).Take(&userConsent); result.Error != nil {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
		return
	}

	if userConsent.ArtefactStatus != models.ArtefactStatusPending {
		// Processing is already done, let's just return.
		return
	}

	consentNotificationMock := consentNotifierStatus{
		ConsentID:     uuid.New(),
		ConsentHandle: userConsent.ConsentHandle,
		ConsentStatus: "ACTIVE",
	}

	if err := updateUserConsentForStatusChange(consentNotificationMock); err != nil {
		handleNotificationError(response, err)
		return
	}

	var updatedUserConsent models.UserConsents
	if result := config.Database.Model(&models.UserConsents{}).Where("user_id = ?",
		userID).Take(&updatedUserConsent); result.Error != nil {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
		return
	}

	config.Database.Where("user_id = ?", userID).Take(&updatedUserConsent)
	if err := saveFipData(updatedUserConsent.SessionID); err != nil {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
		return
	}
}

// Data flow

// FINotification returns notification JSON response
func FINotification(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(request.Body)
	var requestJSON fINotificationRequest
	if err := decoder.Decode(&requestJSON); err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	if err := saveFipData(requestJSON.FIStatusNotification.SessionID); err != nil {
		handleNotificationError(response, err)
		return
	}

	clientAPI, requestJWS, setuResponseBody, err := sendResponseToSetuNotification()
	if err != nil {
		handleNotificationError(response, err)
		return
	}

	response.Header().Set("client_api_key", clientAPI)
	response.Header().Set("x-jws-signature", requestJWS)
	response.Write(setuResponseBody)
}
