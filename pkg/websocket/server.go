package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	filelogger "bitbucket.org/lpi-tech-dev/websocket-backend/pkg/lib/log"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	Clients    map[*Client]bool `json:"clients"`
	Register   chan *Client     `json:"-"`
	Unregister chan *Client     `json:"-"`
	Broadcast  chan []byte      `json:"-"`
}

func NewServer() *Server {
	return &Server{
		Clients:    map[*Client]bool{},
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

// run server
func (server *Server) Run() {
	for {
		select {
		case client := <-server.Register:
			server.RegisterClient(client)

		case client := <-server.Unregister:
			server.UnregisterClient(client)

		case message := <-server.Broadcast:
			server.BroadcastMessage(message)
		}
	}
}

func (server *Server) RegisterClient(client *Client) {
	log.Printf("New client: %s@%s", client.Username, client.Hub)
	server.Clients[client] = true
}

func (server *Server) UnregisterClient(client *Client) {
	if _, ok := server.Clients[client]; ok {
		delete(server.Clients, client)
	}
}

func (server *Server) BroadcastMessage(message []byte) {
	// Step 1
	// before broadcasting unmarshal the message package to get structured data
	m := Message{}
	if err := json.Unmarshal(message, &m); err != nil {
		logMessage := fmt.Sprintf("error structuring message: %v\nerror at: websocket/server.go@65\n", err.Error())
		log.Println(logMessage)
		filelogger.Log(logrus.ErrorLevel, logMessage)
	}

	// Step 2
	// get client on same hub by calling fn GetClientOnHub and user m.Client.Target as its parameter
	clientOnHub := server.GetClientOnHub(m.Client.Hub)

	// Step 3
	// Start broadcast
	for client := range clientOnHub {
		client.Send <- message
	}
}

func (server *Server) GetClientOnHub(hub string) map[*Client]bool {
	onhub := map[*Client]bool{}

	for client := range server.Clients {
		if client.Hub == hub {
			onhub[client] = true
		}
	}

	return onhub
}
