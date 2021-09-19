package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"insurance-ng/src/server/config"
	"insurance-ng/src/server/controllers/account_aggregator"
	"insurance-ng/src/server/controllers/users"
	"insurance-ng/src/server/middleware"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	muxDispatcher := mux.NewRouter()
	config.LoadEnv()
	config.ConnectToDb()
	muxDispatcher.Handle(users.UrlRegister, jwtAuth(users.RegisterUserHandler)).Methods("OPTIONS", "GET")
	// Account Aggregator
	//// Consent Flow
	muxDispatcher.Handle(account_aggregator.UrlConsentStatus,
		jwtAuth(account_aggregator.GetConsentStatus)).Methods("OPTIONS", "GET")
	muxDispatcher.Handle(account_aggregator.UrlCreateConsent,
		jwtAuth(account_aggregator.CreateConsentRequest)).Methods("OPTIONS", "POST")
	//// Data Flow
	muxDispatcher.Handle(account_aggregator.UrlGetUserData,
		jwtAuth(account_aggregator.GetUserData)).Methods("OPTIONS", "GET")
	//// Notifications Flow
	muxDispatcher.Handle(account_aggregator.UrlConsentNotification,
		setuAuth(account_aggregator.ConsentNotification)).Methods("POST")
	// Static Files
	muxDispatcher.PathPrefix("/").Handler(http.FileServer(http.Dir("./dist/")))
	// Start Server
	port := os.Getenv("APP_PORT")
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), muxDispatcher)
	log.Fatal(err)
	fmt.Printf("Mux HTTP server running on port %s", port)
}

func jwtAuth(controller func(http.ResponseWriter, *http.Request)) *negroni.Negroni {
	return negroni.New(
		middleware.CorsMiddleware(),
		negroni.HandlerFunc(middleware.JwtMiddleware().HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(controller)),
	)
}

func setuAuth(controller func(http.ResponseWriter, *http.Request)) *negroni.Negroni {
	return negroni.New(
		negroni.HandlerFunc(middleware.NewSetuValidateMiddleware().HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(controller)),
	)
}
