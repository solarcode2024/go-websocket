package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	filelogger "bitbucket.org/lpi-tech-dev/websocket-backend/pkg/lib/log"
	"github.com/sirupsen/logrus"
)

const (
	BroadcastAction = "Broadcast"
	LeaveAction     = "Leave"
	JoinAction      = "Join"
)

type RawMessage struct {
	Message string `json:"message"`
}

type Message struct {
	Message string    `json:"message"`
	Client  *Client   `json:"sender"`
	Target  string    `json:"target"`
	Action  string    `json:"action"`
	IsMeta  bool      `json:"is_meta"`
	Date    time.Time `json:"date"`
}

func NewMessage(message []byte, action string, client *Client) *Message {
	var isMeta bool
	rm := RawMessage{}
	if err := json.Unmarshal(message, &rm); err != nil {
		logMessage := fmt.Sprintf("error structuring message: %v\nerror at: websocket/message.go@35\n", err.Error())
		log.Println(logMessage)
		filelogger.Log(logrus.ErrorLevel, logMessage)
	}

	switch action {
	case BroadcastAction:
		isMeta = false
	default:
		isMeta = true
	}

	return &Message{
		Message: rm.Message,
		Client:  client,
		Target:  client.Hub,
		Action:  action,
		IsMeta:  isMeta,
		Date:    time.Now(),
	}
}

func (message *Message) Marhsal() ([]byte, error) {
	marshaled, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return marshaled, nil
}
