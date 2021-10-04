package insurance

import (
	"encoding/json"
	"insurance-ng/src/server/controllers"
	"net/http"
)

const (
	UrlInsurancePurchase = "/api/v1/insurance/purchase/"
	UrlInsuranceClaim    = "/api/v1/insurance/claim/"
	UrlGetUserData       = "/api/v1/insurance/"
	UrlWaitForProcessing = "/api/v1/ws/insurance/"
)

func GetUserData(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userId, err := controllers.GetUserIdentifier(response, request)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	userData, err := getUserDataRecord(userId)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	respMessage, _ := json.Marshal(userData)
	response.Write(respMessage)
}

func InsurancePurchaseHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userId, err := controllers.GetUserIdentifier(response, request)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	insuranceUuid, err := getInsuranceUuid(request.Body)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	if err := createInsuranceRecord(userId, insuranceUuid); err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	userData, err := getUserDataRecord(userId)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	respMessage, _ := json.Marshal(userData)
	response.Write(respMessage)
}

func InsuranceClaimHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userId, err := controllers.GetUserIdentifier(response, request)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	insuranceUuid, err := getInsuranceUuid(request.Body)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	if err := initiateInsuranceClaim(userId, insuranceUuid); err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	userData, err := getUserDataRecord(userId)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	respMessage, _ := json.Marshal(userData)
	response.Write(respMessage)
}

func WaitForDataProcessingWebsocket(websocket Websocket, response http.ResponseWriter, request *http.Request) {
	connection, err := wsUpgrader.Upgrade(response, request, nil)
	if err != nil {
		return
	}
	client := &Client{websocket: &websocket, connection: connection}
	client.websocket.register <- client

	go client.websocketDataFetchedSignal()
}
