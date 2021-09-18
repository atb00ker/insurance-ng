package account_aggregator

import (
	"encoding/json"
	"insurance-ng/src/server/controllers"
	"net/http"
)

const (
	UrlCreateConsent = "/api/account_aggregator/consent/"
	UrlConsentStatus = "/api/account_aggregator/consent/status/"
	UrlGetUserData   = "/api/account_aggregator/data/"
)

func CreateConsentRequest(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userId, ok := controllers.GetUserIdentifier(response, request)
	if !ok {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
		return
	}

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	var requestJson createConsentRequestData
	if err := decoder.Decode(&requestJson); err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	consentResponse, expiry, err := sendCreateConsentReqToAcctAggregator(requestJson.Phone)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}
	result := addOrUpdateConsentToDb(userId, consentResponse, expiry)

	if result.Error != nil {
		controllers.HandleError(response, result.Error.Error())
		return
	}
	respMessage, _ := json.Marshal(createConsentResponseData{
		ConsentHandle: consentResponse.ConsentHandle.String(),
	})
	response.Write(respMessage)
}

func GetConsentStatus(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userId, ok := controllers.GetUserIdentifier(response, request)
	if !ok {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
		return
	}

	userConsent, err := getUserConsent(userId)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	if userConsent.UserData != "-" {
		respMessage, _ := json.Marshal(controllers.ResponseMessage{Status: userConsent.Status})
		response.Write(respMessage)
		return
	}

	consentId, err := getUserArtefactStatus(userId, userConsent.ConsentHandle)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	status, err := fetchSignedConsent(userId, consentId)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	respMessage, _ := json.Marshal(controllers.ResponseMessage{Status: status})
	response.Write(respMessage)
}

func GetUserData(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userId, ok := controllers.GetUserIdentifier(response, request)
	if !ok {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
		return
	}

	userConsent, err := getUserConsent(userId)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	var userData []rahasyaDataResponseCollection
	rahasyaKeys := getRahasyaKeys()
	if userConsent.UserData == "-" {
		sessionData, err := getDataSession(userId, rahasyaKeys, userConsent)
		if err != nil {
			controllers.HandleError(response, err.Error())
			return
		}

		encryptedData, err := getEncryptedFIData(sessionData)
		if err != nil {
			controllers.HandleError(response, err.Error())
			return
		}

		userData, err = getUnencryptedFIDataList(rahasyaKeys, encryptedData)
		if err != nil {
			controllers.HandleError(response, err.Error())
			return
		}

		// TODO: Bad Idea to store unencrypted data
		if err := savebase64Data(userData, userId); err != nil {
			controllers.HandleError(response, err.Error())
			return
		}
	} else {
		if err := json.Unmarshal([]byte(userConsent.UserData), &userData); err != nil {
			controllers.HandleError(response, err.Error())
			return
		}
	}

	respMessage, _ := json.Marshal(userData)
	response.Write(respMessage)
}
