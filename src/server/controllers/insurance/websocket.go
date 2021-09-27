package insurance

import (
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
		// TODO
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
			return
		}

		response, err := client.connection.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}

		response.Write([]byte(strconv.FormatBool(userConsent.DataFetched)))
		if err := response.Close(); err != nil {
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

			response, err := client.connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			response.Write([]byte(strconv.FormatBool(true)))
			if err := response.Close(); err != nil {
				return
			}
		}
	}
}
