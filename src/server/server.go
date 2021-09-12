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
	muxDispatcher.Handle(users.UrlRegister, authenticated(users.RegisterUserHandler)).Methods("OPTIONS", "GET")
	muxDispatcher.Handle(account_aggregator.UrlConsentStatus,
		authenticated(account_aggregator.GetConsentStatus)).Methods("OPTIONS", "GET")
	muxDispatcher.Handle(account_aggregator.UrlGetUserData,
		authenticated(account_aggregator.GetUserData)).Methods("OPTIONS", "GET")
	muxDispatcher.Handle(account_aggregator.UrlCreateConsent,
		authenticated(account_aggregator.CreateConsentRequest)).Methods("OPTIONS", "POST")
	// Static Files
	muxDispatcher.PathPrefix("/").Handler(http.FileServer(http.Dir("./dist/")))
	// Start Server
	port := os.Getenv("APP_PORT")
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), muxDispatcher)
	log.Fatal(err)
	fmt.Printf("Mux HTTP server running on port %s", port)
}

func authenticated(controller func(http.ResponseWriter, *http.Request)) *negroni.Negroni {
	return negroni.New(
		middleware.CorsMiddleware(),
		negroni.HandlerFunc(middleware.JwtMiddleware().HandlerWithNext),
		negroni.Wrap(http.HandlerFunc(controller)),
	)
}
