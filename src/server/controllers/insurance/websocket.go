package insurance

import (
	"errors"
	"fmt"
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

// WaitForProcessing  signal that will be triggered
// when a new FIP data is processed.
var WaitForProcessing = make(chan string)

const (
	dataNotShared     string = "data-not-shared"
	consentNotStarted        = "consent-not-started"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Limit origin of websocket.
		return true
	},
}

// NewWebocket allows server to create a new websocket server instance
func NewWebocket() *Websocket {
	// Websocket instance containing the state information of
	// the socket.
	return &Websocket{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*string]*Client),
	}
}

// Init infinitely checks for any messages for websocket
// and performs appropriate actions
func (socket *Websocket) Init() {
	for {
		select {
		case client := <-socket.register:
			socket.clients[client.ID] = client
		case client := <-socket.unregister:
			delete(socket.clients, client.ID)
		}
	}
}

func (client *Client) websocketDataFetchedSignal() {
	defer func() {
		client.websocket.unregister <- client
		client.connection.Close()
	}()

	for {
		_, message, err := client.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		userID := string(message)
		userConsent, err := getUserConsentWithUserID(userID)
		if err != nil {
			if err := websocketResponse(client, []byte(err.Error())); err != nil {
				return
			}
			continue
		}

		// If the data is ready, send the signal for the same.
		if err := websocketResponse(client, []byte(strconv.FormatBool(userConsent.DataFetched))); err != nil {
			return
		}

		if !userConsent.DataFetched {
			var completedProcessUserID string
			for {
				completedProcessUserID = <-WaitForProcessing
				if completedProcessUserID == userID {
					break
				}
			}

			if checkUserScore(userConsent) == 0 {
				// Not enough data sources are shared.
				if err = websocketResponse(client, []byte(dataNotShared)); err != nil {
					return
				}
				continue
			}

			// We are ready to show data.
			if err := websocketResponse(client, []byte(strconv.FormatBool(true))); err != nil {
				return
			}
		}
	}
}

func getUserConsentWithUserID(userID string) (userConsent models.UserConsents, err error) {
	if result := config.Database.Model(&models.UserConsents{}).Where("user_id = ?",
		userID).Take(&userConsent); result.Error != nil {
		if result.Error.Error() == "record not found" {
			err = errors.New(consentNotStarted)
			return
		}
	}

	if checkUserScore(userConsent) == 0 {
		err = errors.New(dataNotShared)
		return
	}

	return userConsent, nil
}

func checkUserScore(userConsent models.UserConsents) int16 {
	var userScores models.UserScores
	if result := config.Database.Model(&models.UserScores{}).Where("user_consent_id = ?",
		userConsent.ID).Take(&userScores); result.Error != nil {
		return -1
	}

	fmt.Println(userScores.SharedDataSources)
	return userScores.SharedDataSources
}
