package users

import (
	"encoding/json"
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/controllers"
	"insurance-ng/src/server/models"
	"net/http"
)

const (
	UrlRegister = "/api/v1/register/"
)

func RegisterUserHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userId, err := controllers.GetUserIdentifier(response, request)
	if err != nil {
		controllers.HandleError(response, err.Error())
		return
	}

	result := config.Database.FirstOrCreate(&models.Users{Id: userId, IsAdmin: false})
	if result.Error != nil {
		controllers.HandleError(response, result.Error.Error())
		return
	}

	respMessage, _ := json.Marshal(controllers.ResponseMessage{Status: controllers.ResponseSuccess})
	response.Write(respMessage)
}
