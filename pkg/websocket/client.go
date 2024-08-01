package websocket

import (
	"fmt"
	"log"
	"time"

	filelogger "bitbucket.org/lpi-tech-dev/websocket-backend/pkg/lib/log"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	writeWait      = 10 * time.Second    // Max wait time when writing message to peer
	pongWait       = 60 * time.Second    // Max time till next pong from peer
	pingPeriod     = (pongWait * 9) / 10 // Send ping interval, must be less then pong wait time
	maxMessageSize = 10000               // Maximum message size allowed from peer.
)

type Client struct {
	ID         uuid.UUID       `json:"id"`
	Username   string          `json:"username"`
	Hub        string          `json:"hub"`
	Connection *websocket.Conn `json:"-"`
	Server     *Server         `json:"-"`
	Send       chan []byte     `json:"-"`
}

func NewClient(username string, hub string, connection *websocket.Conn, server *Server) *Client {
	return &Client{
		ID:         uuid.New(),
		Username:   username,
		Hub:        hub,
		Connection: connection,
		Server:     server,
		Send:       make(chan []byte),
	}
}

// to read message that send over the websocket
// or to read message sended by user
func (client *Client) ReadEngine() {
	defer func() {
		client.Disconnect()
	}()

	client.Connection.SetReadLimit(maxMessageSize)
	client.Connection.SetReadDeadline(time.Now().Add(pongWait))
	client.Connection.SetPongHandler(func(string) error { client.Connection.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := client.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
				logMessage := fmt.Sprintf("error read message: %v\nerror at: websocket/client.go@55\n", err.Error())
				log.Println(logMessage)
				filelogger.Log(logrus.ErrorLevel, logMessage)
			}
			break
		}

		MessagePackage := NewMessage(message, BroadcastAction, client)
		MarshaledMessagePackage, err := MessagePackage.Marhsal()
		if err != nil {
			logMessage := fmt.Sprintf("error marshaling message: %v\nerror at: websocket/client.go@64\n", err.Error())
			log.Println(logMessage)
			filelogger.Log(logrus.ErrorLevel, logMessage)
		}
		client.Server.Broadcast <- MarshaledMessagePackage
	}
}

func (client *Client) WriteEngine() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.Connection.Close()
	}()

	for {
		select {
		case message, ok := <-client.Send:
			client.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.Connection.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-client.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}

func (client *Client) Disconnect() {
	client.Server.Unregister <- client
	close(client.Send)
	client.Connection.Close()
}
