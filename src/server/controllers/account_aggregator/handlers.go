package account_aggregator

import (
	"encoding/json"
	"fmt"
	"insurance-ng/src/server/controllers"
	"net/http"
	"time"
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
	startTime := time.Now()
	startTimeHack := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.153Z",
		startTime.Year(), startTime.Month(), startTime.Day(),
		startTime.Hour(), startTime.Minute(), startTime.Second())

	if err := updateUserConsentForStatusChange(requestJson.ConsentStatusNotification); err != nil {
		HandleNotificationError(response, startTimeHack, err.Error())
		return
	}

	clientApi, requestJws, setuResponseBody, err := sendResponseToSetuNotification(startTimeHack)
	if err != nil {
		HandleNotificationError(response, startTimeHack, err.Error())
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
	startTime := time.Now()
	startTimeHack := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.153Z",
		startTime.Year(), startTime.Month(), startTime.Day(),
		startTime.Hour(), startTime.Minute(), startTime.Second())

	if err := saveFipData(requestJson.FIStatusNotification.SessionID); err != nil {
		HandleNotificationError(response, startTimeHack, err.Error())
		return
	}

	clientApi, requestJws, setuResponseBody, err := sendResponseToSetuNotification(startTimeHack)
	if err != nil {
		HandleNotificationError(response, startTimeHack, err.Error())
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
