package middleware

import (
	"insurance-ng/src/server/controllers"
	"net/http"
	"os"
)

type SetuValidateMiddleware struct{}

func NewSetuValidateMiddleware() *SetuValidateMiddleware {
	return &SetuValidateMiddleware{}
}

func (m *SetuValidateMiddleware) HandlerWithNext(response http.ResponseWriter, request *http.Request,
	next http.HandlerFunc) {
	auth := request.Header.Get("Authorization")
	if auth == os.Getenv("APP_SETU_AA_KEY") {
		// TODO: Check JWS
		next(response, request)
	} else {
		controllers.HandleError(response, "Incorrect Authorization token!")
	}
}
