package controllers

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/form3tech-oss/jwt-go"
)

const DataEndpointsExplainationInParagraph = `To provide the lowest prices, we need your financial information to understand you. We take the data points like personal information (including name, date of birth, address, pancard & existing insurance plans), deposit account transactions, current balance summary, insurance accounts and transactions, investment plans and debt to understand your lifestyle and plans that would best suit you. You have the right to revoke/request for data deletion at any point in the future. Sharing bank account and insurance account is a minimum requirement.`

func HandleError(response http.ResponseWriter, err string) {
	response.WriteHeader(http.StatusBadRequest)
	respMessage, _ := json.Marshal(ResponseMessage{
		Status: ResponseError,
		Error:  err,
	})
	response.Write(respMessage)
}

func GetUserIdentifier(response http.ResponseWriter, request *http.Request) (string, error) {
	token := request.Context().Value("user")
	user := token.(*jwt.Token).Claims.(jwt.MapClaims)
	userId, ok := user["sub"].(string)
	if !ok {
		return "", errors.New(IsUserLoggedInErrorMessage)
	}
	return userId, nil
}

func GetRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	account_id := make([]rune, length)
	for i := range account_id {
		account_id[i] = letters[rand.Intn(len(letters))]
	}
	return string(account_id)
}
