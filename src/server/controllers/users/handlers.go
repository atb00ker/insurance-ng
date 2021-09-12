package users

import (
	"encoding/json"
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/controllers"
	"insurance-ng/src/server/models"
	"net/http"
)

const (
	UrlRegister = "/register/"
)

func RegisterUserHandler(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	userId, ok := controllers.GetUserIdentifier(response, request)
	if !ok {
		controllers.HandleError(response, controllers.IsUserLoggedInErrorMessage)
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
