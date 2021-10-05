package controllers

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/form3tech-oss/jwt-go"
)

// HandleError creates the json send in error message and sends
func HandleError(response http.ResponseWriter, err string) {
	response.WriteHeader(http.StatusBadRequest)
	respMessage, _ := json.Marshal(ResponseMessage{
		Status: ResponseError,
		Error:  err,
	})
	response.Write(respMessage)
}

// GetUserIdentifier gets user from the parsed JWT body
func GetUserIdentifier(response http.ResponseWriter, request *http.Request) (string, error) {
	token := request.Context().Value("user")
	user := token.(*jwt.Token).Claims.(jwt.MapClaims)
	userID, ok := user["sub"].(string)
	if !ok {
		return "", errors.New(IsUserLoggedInErrorMessage)
	}
	return userID, nil
}

// GetRandomString generates a random string in range [A-Z,0-9]
func GetRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	accountID := make([]rune, length)
	for i := range accountID {
		accountID[i] = letters[rand.Intn(len(letters))]
	}
	return string(accountID)
}
