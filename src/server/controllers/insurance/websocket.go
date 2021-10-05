package insurance

import (
	"fmt"
	"insurance-ng/src/server/config"
	"insurance-ng/src/server/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var WaitForProcessing = make(chan string)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Limit origin of websocket.
		return true
	},
}

func NewWebocket() *Websocket {
	return &Websocket{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*string]*Client),
	}
}

func (socket *Websocket) Init() {
	for {
		select {
		case client := <-socket.register:
			socket.clients[client.Id] = client
		case client := <-socket.unregister:
			delete(socket.clients, client.Id)
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

		userId := string(message)
		var userConsent models.UserConsents
		if result := config.Database.Model(&models.UserConsents{}).Where("user_id = ?",
			userId).Take(&userConsent); result.Error != nil {

			if result.Error.Error() == "record not found" {
				if err := websocketResponse(client, []byte("consent-not-started")); err != nil {
					return
				}
				// If the consent is not even started, we don't
				// need to wait for the next signal.
				continue
			}
			return
		}

		if client.checkUserScore(userConsent) == 0 {
			// Not enough data sources are shared.
			if err = websocketResponse(client, []byte("data-not-shared")); err != nil {
				return
			}
			continue
		}

		// If the data is ready, send the signal for the same.
		if err := websocketResponse(client, []byte(strconv.FormatBool(userConsent.DataFetched))); err != nil {
			return
		}

		if !userConsent.DataFetched {
			var completedProcessUserId string
			for {
				completedProcessUserId = <-WaitForProcessing
				if completedProcessUserId == userId {
					break
				}
			}

			if client.checkUserScore(userConsent) == 0 {
				// Not enough data sources are shared.
				if err = websocketResponse(client, []byte("data-not-shared")); err != nil {
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

func (client *Client) checkUserScore(userConsent models.UserConsents) int16 {
	var userScores models.UserScores
	if result := config.Database.Model(&models.UserScores{}).Where("user_consent_id = ?",
		userConsent.Id).Take(&userScores); result.Error != nil {
		return -1
	}

	fmt.Println(userScores.SharedDataSources)
	return userScores.SharedDataSources
}
