package users

import (
	"encoding/json"
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/controllers"
	"insurance-ng/src/server/models"
	"net/http"
)

// URLs for user endpoints
const (
	URLRegister = "/api/v1/register/"
)

// RegisterUserHandler registers user on their first login
func RegisterUserHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userID, err := controllers.GetUserIdentifier(response, request)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	result := config.Database.FirstOrCreate(&models.Users{ID: userID, IsAdmin: false})
	if result.Error != nil {
		controllers.HandleError(response, result.Error.Error())
		return
	}

	respMessage, _ := json.Marshal(controllers.ResponseMessage{Status: controllers.ResponseSuccess})
	response.Write(respMessage)
}
