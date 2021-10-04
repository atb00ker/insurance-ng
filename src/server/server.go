package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"insurance-ng/src/server/config"
	"insurance-ng/src/server/controllers/account_aggregator"
	"insurance-ng/src/server/controllers/insurance"
	"insurance-ng/src/server/controllers/users"
	"insurance-ng/src/server/middleware"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	websocket := insurance.NewWebocket()
	go websocket.Init()
	muxDispatcher := mux.NewRouter()
	config.LoadEnv()
	config.ConnectToDb()
	// User
	muxDispatcher.Handle(users.UrlRegister,
		jwtAuth(users.RegisterUserHandler)).Methods("OPTIONS", "GET")
	// Insurance
	muxDispatcher.Handle(insurance.UrlInsurancePurchase,
		jwtAuth(insurance.InsurancePurchaseHandler)).Methods("OPTIONS", "POST")
	muxDispatcher.Handle(insurance.UrlInsuranceClaim,
		jwtAuth(insurance.InsuranceClaimHandler)).Methods("OPTIONS", "POST")
	muxDispatcher.Handle(insurance.UrlGetUserData,
		jwtAuth(insurance.GetUserData)).Methods("OPTIONS", "GET")
	muxDispatcher.Handle(insurance.UrlWaitForProcessing,
		insuranceWebsocket(insurance.WaitForDataProcessingWebsocket, websocket))
	// Account Aggregator
	muxDispatcher.Handle(account_aggregator.UrlCreateConsent,
		jwtAuth(account_aggregator.CreateConsentRequest)).Methods("OPTIONS", "POST")
	muxDispatcher.Handle(account_aggregator.UrlConsentNotificationMock,
		jwtAuth(account_aggregator.ConsentNotificationMock)).Methods("OPTIONS", "GET")
	muxDispatcher.Handle(account_aggregator.UrlConsentNotification,
		setuAuth(account_aggregator.ConsentNotification)).Methods("POST")
	muxDispatcher.Handle(account_aggregator.UrlArtefactNotification,
		setuAuth(account_aggregator.FINotification)).Methods("POST")
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

func insuranceWebsocket(controller func(websocket insurance.Websocket, response http.ResponseWriter,
	request *http.Request), websocket *insurance.Websocket) *negroni.Negroni {
	return negroni.New(
		middleware.CorsMiddleware(),
		negroni.Wrap(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			controller(*websocket, response, request)
		})),
	)
}
