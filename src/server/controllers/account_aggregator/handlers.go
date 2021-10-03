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
	UrlCreateConsent           = "/api/account_aggregator/consent/"
	UrlConsentNotificationMock = "/api/account_aggregator/Mock/Consent/Notification/"
	UrlConsentNotification     = "/api/account_aggregator/Consent/Notification"
	UrlArtefactNotification    = "/api/account_aggregator/FI/Notification"
)

func CreateConsentRequest(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userId, ok := controllers.GetUserIdentifier(response, request)
	if !ok {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
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

	if result := addOrUpdateConsentToDb(userId, consentResponse, expiry); result.Error != nil {
		controllers.HandleError(response, result.Error.Error())
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
	// This is just a mock endpoint.
	// We don't need in it real life scenarios.
	response.Header().Set("Content-Type", "application/json")
	userId, ok := controllers.GetUserIdentifier(response, request)
	if !ok {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
		return
	}

	var userConsent models.UserConsents
	if result := config.Database.Model(&models.UserConsents{}).Where("user_id = ?",
		userId).Take(&userConsent); result.Error != nil {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
		return
	}

	if userConsent.ArtefactStatus == models.ArtefactStatusPending {
		consentNotificationMock := consentNotifierStatus{
			ConsentId:     uuid.New(),
			ConsentHandle: userConsent.ConsentHandle,
			ConsentStatus: "ACTIVE",
		}
		if err := updateUserConsentForStatusChange(consentNotificationMock); err != nil {
			HandleNotificationError(response, err)
			return
		}
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

	// Hack:
	// TODO - Currently the Setu FI notification is not triggered
	// sometimes, which can be a big problem in the
	// hackathon, hence, for the time being, I am commenting this
	// out and focusing and trigger this notification automatically
	// after consent request.
	// if err := saveFipData(requestJson.FIStatusNotification.SessionID); err != nil {
	// 	HandleNotificationError(response, err)
	// 	return
	// }

	clientApi, requestJws, setuResponseBody, err := sendResponseToSetuNotification()
	if err != nil {
		HandleNotificationError(response, err)
		return
	}

	response.Header().Set("client_api_key", clientApi)
	response.Header().Set("x-jws-signature", requestJws)
	response.Write(setuResponseBody)
}
