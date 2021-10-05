package middleware

import (
	"insurance-ng/src/server/controllers"
	"net/http"
	"os"
)

// SetuValidateMiddleware is a type for server to create middleware object
type SetuValidateMiddleware struct{}

// NewSetuValidateMiddleware initiates new setu validation object
func NewSetuValidateMiddleware() *SetuValidateMiddleware {
	return &SetuValidateMiddleware{}
}

// HandlerWithNext validates and call next on success
func (m *SetuValidateMiddleware) HandlerWithNext(response http.ResponseWriter, request *http.Request,
	next http.HandlerFunc) {
	auth := request.Header.Get("Authorization")
	if auth == os.Getenv("APP_SETU_AA_KEY") {
		// TODO: Check Setu JWS
		next(response, request)
	} else {
		controllers.HandleError(response, "Incorrect Authorization token!")
	}
}
