package websocket

import (
	"testing"

	"github.com/google/uuid"
)

func TestGetClientOnHub(t *testing.T) {
	client1 := Client{
		ID:       uuid.New(),
		Username: "astra",
		Hub:      "test-hub",
	}
	client2 := Client{
		ID:       uuid.New(),
		Username: "nova",
		Hub:      "test-hub",
	}
	client3 := Client{
		ID:       uuid.New(),
		Username: "hydra",
		Hub:      "test-hub-unified",
	}

	clients := map[*Client]bool{
		&client1: true,
		&client2: true,
		&client3: true,
	}
	server := Server{
		Clients:    clients,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
	clientOnHub := server.GetClientOnHub("test-hub")
	if len(clientOnHub) == 0 {
		t.Fatalf("no client on hub")
	}

	for v := range clientOnHub {
		t.Log(v)
	}
}
