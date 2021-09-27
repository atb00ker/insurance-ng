package insurance

import (
	"encoding/json"
	"insurance-ng/src/server/controllers"
	"net/http"
)

const (
	InsurancePurchase    = "/api/insurance/purchase/"
	UrlGetUserData       = "/api/insurance/"
	UrlWaitForProcessing = "/api/ws/insurance/"
)

func GetUserData(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userId, ok := controllers.GetUserIdentifier(response, request)
	if !ok {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
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
	userId, ok := controllers.GetUserIdentifier(response, request)
	if !ok {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
		return
	}

	decoder := json.NewDecoder(request.Body)
	var requestJson purchaseRequest
	if err := decoder.Decode(&requestJson); err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	if err := createInsuranceRecord(userId, requestJson.Uuid); err != nil {
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
