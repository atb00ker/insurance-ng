package controllers

import (
	"encoding/json"
	"math/rand"
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

func GetRandomString(n int) string {
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
