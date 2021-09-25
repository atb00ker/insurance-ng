package account_aggregator

import (
	"encoding/json"
	"insurance-ng/src/server/controllers"
	"net/http"
)

const (
	UrlCreateConsent        = "/api/account_aggregator/consent/"
	UrlConsentStatus        = "/api/account_aggregator/consent/status/"
	UrlGetUserData          = "/api/account_aggregator/data/"
	UrlConsentNotification  = "/api/account_aggregator/Consent/Notification"
	UrlArtefactNotification = "/api/account_aggregator/FI/Notification"
)

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

func GetUserData(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userId, ok := controllers.GetUserIdentifier(response, request)
	if !ok {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
		return
	}

	userData, err := getUserData(userId)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	respMessage, _ := json.Marshal(userData)
	response.Write(respMessage)
}
