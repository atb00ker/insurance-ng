package insurance

import (
	"encoding/json"
	"insurance-ng/src/server/controllers"
	"net/http"
)

// URLs for insurance endpoints
const (
	URLInsurancePurchase = "/api/v1/insurance/purchase/"
	URLInsuranceClaim    = "/api/v1/insurance/claim/"
	URLGetUserData       = "/api/v1/insurance/"
	URLWaitForProcessing = "/api/v1/ws/insurance/"
)

// GetUserData gets user's data procured by FIP
func GetUserData(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userID, err := controllers.GetUserIdentifier(response, request)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	userData, err := getUserDataRecord(userID)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	respMessage, _ := json.Marshal(userData)
	response.Write(respMessage)
}

// InitiatePurchase handles insurance purchase requests
func InitiatePurchase(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userID, err := controllers.GetUserIdentifier(response, request)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	insuranceUUID, err := getInsuranceUUID(request.Body)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	if err := createInsuranceRecord(userID, insuranceUUID); err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	userData, err := getUserDataRecord(userID)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	respMessage, _ := json.Marshal(userData)
	response.Write(respMessage)
}

// InitiateClaim handles insurance claim requests
func InitiateClaim(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userID, err := controllers.GetUserIdentifier(response, request)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	insuranceUUID, err := getInsuranceUUID(request.Body)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	if err := initiateInsuranceClaim(userID, insuranceUUID); err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	userData, err := getUserDataRecord(userID)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	respMessage, _ := json.Marshal(userData)
	response.Write(respMessage)
}

// WaitForDataProcessingWebsocket starts a go routine to open a websocket with client
func WaitForDataProcessingWebsocket(websocket Websocket, response http.ResponseWriter, request *http.Request) {
	connection, err := wsUpgrader.Upgrade(response, request, nil)
	if err != nil {
		return
	}
	client := &Client{websocket: &websocket, connection: connection}
	client.websocket.register <- client

	go client.websocketDataFetchedSignal()
}
