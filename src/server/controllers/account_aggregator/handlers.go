package account_aggregator

import (
	"encoding/json"
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/controllers"
	"insurance-ng/src/server/models"
	"net/http"

	"github.com/google/uuid"
)

const (
	UrlCreateConsent           = "/api/v1/account_aggregator/consent/"
	UrlConsentNotificationMock = "/api/v1/account_aggregator/Mock/Consent/Notification/"
	UrlConsentNotification     = "/api/v1/account_aggregator/Consent/Notification"
	UrlArtefactNotification    = "/api/v1/account_aggregator/FI/Notification"
)

func CreateConsentRequest(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userId, err := controllers.GetUserIdentifier(response, request)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	decoder := json.NewDecoder(request.Body)
	var requestJson createConsentRequestInput
	if err := decoder.Decode(&requestJson); err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	consentResponse, expiry, err := sendCreateConsentReqToAcctAggregator(requestJson.Phone)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	if _, err = addOrUpdateConsentToDb(userId, consentResponse, expiry); err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	respMessage, _ := json.Marshal(createConsentResponseOutput{
		ConsentHandle: consentResponse.ConsentHandle.String(),
	})

	response.Write(respMessage)
}

func ConsentNotification(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(request.Body)
	var requestJson setuConsentNotificationRequest
	if err := decoder.Decode(&requestJson); err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	if err := updateUserConsentForStatusChange(requestJson.ConsentStatusNotification); err != nil {
		HandleNotificationError(response, err)
		return
	}

	clientApi, requestJws, setuResponseBody, err := sendResponseToSetuNotification()
	if err != nil {
		HandleNotificationError(response, err)
		return
	}

	response.Header().Set("client_api_key", clientApi)
	response.Header().Set("x-jws-signature", requestJws)
	response.Write(setuResponseBody)
}

func ConsentNotificationMock(response http.ResponseWriter, request *http.Request) {
	// Notification Hack:
	// This is just a mock endpoint.
	// Currently the Setu notifications are not triggered
	// sometimes, which can be a big problem in the
	// hackathon, hence, for the time being, we start the
	// notification steps here. Not required in production.
	userId, err := controllers.GetUserIdentifier(response, request)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	var userConsent models.UserConsents
	if result := config.Database.Model(&models.UserConsents{}).Where("user_id = ?",
		userId).Take(&userConsent); result.Error != nil {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
		return
	}

	if userConsent.ArtefactStatus != models.ArtefactStatusPending {
		// Processing is already done, let's just return.
		return
	}

	consentNotificationMock := consentNotifierStatus{
		ConsentId:     uuid.New(),
		ConsentHandle: userConsent.ConsentHandle,
		ConsentStatus: "ACTIVE",
	}

	if err := updateUserConsentForStatusChange(consentNotificationMock); err != nil {
		HandleNotificationError(response, err)
		return
	}

	var updatedUserConsent models.UserConsents
	if result := config.Database.Model(&models.UserConsents{}).Where("user_id = ?",
		userId).Take(&updatedUserConsent); result.Error != nil {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
		return
	}

	config.Database.Where("user_id = ?", userId).Take(&updatedUserConsent)
	if err := saveFipData(updatedUserConsent.SessionId); err != nil {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
		return
	}
}

// Data flow
func FINotification(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(request.Body)
	var requestJson fINotificationRequest
	if err := decoder.Decode(&requestJson); err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	if err := saveFipData(requestJson.FIStatusNotification.SessionID); err != nil {
		HandleNotificationError(response, err)
		return
	}

	clientApi, requestJws, setuResponseBody, err := sendResponseToSetuNotification()
	if err != nil {
		HandleNotificationError(response, err)
		return
	}

	response.Header().Set("client_api_key", clientApi)
	response.Header().Set("x-jws-signature", requestJws)
	response.Write(setuResponseBody)
}
