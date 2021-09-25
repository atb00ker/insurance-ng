package insurance

import (
	"encoding/json"
	"fmt"
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/controllers"
	"insurance-ng/src/server/models"
	"net/http"
)

const (
	InsurancePurchase = "/api/insurance/purchase/"
)

func InsurancePurchaseHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userId, ok := controllers.GetUserIdentifier(response, request)
	if !ok {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
		return
	}

	fmt.Println(userId)
	decoder := json.NewDecoder(request.Body)
	var requestJson purchaseRequest
	if err := decoder.Decode(&requestJson); err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	var insurance models.Insurance
	result := config.Database.Model(&models.Insurance{}).Where("id = ?",
		requestJson.Uuid).Take(&insurance)
	if result.Error != nil {
		controllers.HandleError(response, result.Error.Error())
		return
	}

	respMessage, _ := json.Marshal(controllers.ResponseMessage{Status: controllers.ResponseSuccess})
	response.Write(respMessage)
}
