package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	responseWriter "bitbucket.org/lpi-tech-dev/websocket-backend/pkg/lib/response"
	websocket "bitbucket.org/lpi-tech-dev/websocket-backend/pkg/websocket"
)

const (
	version  = "1.0.0"
	port     = ":8001"
	certPath = "cert/cert.pem"
	keyPath  = "cert/key.pem"
)

func main() {
	// initiate websocket server
	WsInstance := websocket.NewServer()
	go WsInstance.Run()

	// monitor handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := responseWriter.MapInterface{
			"message": "hello fellow coders!",
			"data": responseWriter.MapInterface{
				"application information": responseWriter.MapString{
					"version":     version,
					"description": "websocket application",
				},
			},
		}

		responseWriter.ResponseWritter(w, data, 200)
	})

	// websocket server handler
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWebsocket(WsInstance, w, r)
	})

	// check cert files
	_, certErr := os.Stat(certPath)
	_, keyErr := os.Stat(keyPath)
	if errors.Is(certErr, os.ErrNotExist) || errors.Is(keyErr, os.ErrNotExist) {
		log.Printf("server up: http://localhost%v\n", port)
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Printf("error server: %v\n", err.Error())
		}
	}

	// serve https server
	log.Printf("server up: https://localhost%v\n", port)
	err := http.ListenAndServeTLS(port, certPath, keyPath, nil)
	if err != nil {
		log.Printf("error server: %v\n", err.Error())
	}
}
