package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/form3tech-oss/jwt-go"
)

func HandleError(response http.ResponseWriter, err string) {
	response.WriteHeader(http.StatusBadRequest)
	respMessage, _ := json.Marshal(ResponseMessage{
		Status: ResponseError,
		Error:  err,
	})
	response.Write(respMessage)
}

func GetUserIdentifier(response http.ResponseWriter, request *http.Request) (string, bool) {
	token := request.Context().Value("user")
	user := token.(*jwt.Token).Claims.(jwt.MapClaims)
	userId, ok := user["sub"].(string)
	return userId, ok
}
