package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"insurance-ng/src/server/config"
	"insurance-ng/src/server/controllers/accountaggregator"
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
	muxDispatcher.Handle(users.URLRegister,
		jwtAuth(users.RegisterUserHandler)).Methods("OPTIONS", "GET")
	// Insurance
	muxDispatcher.Handle(insurance.URLInsurancePurchase,
		jwtAuth(insurance.InitiatePurchase)).Methods("OPTIONS", "POST")
	muxDispatcher.Handle(insurance.URLInsuranceClaim,
		jwtAuth(insurance.InitiateClaim)).Methods("OPTIONS", "POST")
	muxDispatcher.Handle(insurance.URLGetUserData,
		jwtAuth(insurance.GetUserData)).Methods("OPTIONS", "GET")
	muxDispatcher.Handle(insurance.URLWaitForProcessing,
		insuranceWebsocket(insurance.WaitForDataProcessingWebsocket, websocket))
	// Account Aggregator
	muxDispatcher.Handle(accountaggregator.URLCreateConsent,
		jwtAuth(accountaggregator.CreateConsentRequest)).Methods("OPTIONS", "POST")
	muxDispatcher.Handle(accountaggregator.URLConsentNotificationMock,
		jwtAuth(accountaggregator.ConsentNotificationMock)).Methods("OPTIONS", "GET")
	muxDispatcher.Handle(accountaggregator.URLConsentNotification,
		setuAuth(accountaggregator.ConsentNotification)).Methods("POST")
	muxDispatcher.Handle(accountaggregator.URLArtefactNotification,
		setuAuth(accountaggregator.FINotification)).Methods("POST")
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
