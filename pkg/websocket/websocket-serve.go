package websocket

import (
	"log"
	"net/http"
)

const WS_LOG_FILE = "websocket.txt"

// check origin off client
func CheckOrigin(r *http.Request) bool { return true }

// to serve websocket server
func ServeWebsocket(wsServer *Server, w http.ResponseWriter, r *http.Request) {
	// get username
	username := r.URL.Query().Get("username")
	if username == "" {
		username = "anonymous"
	}

	// get hub
	hub := r.URL.Query().Get("hub")
	if hub == "" {
		hub = "public"
	}

	// upgrade protocol
	Upgrader.CheckOrigin = CheckOrigin
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrade protocol: %v", err)
		return
	}

	// generate client
	client := NewClient(username, hub, conn, wsServer)

	// call write & read engine
	go client.ReadEngine()
	go client.WriteEngine()

	wsServer.Register <- client
}
