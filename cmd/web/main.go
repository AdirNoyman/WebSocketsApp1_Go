package main

import (
	"log"
	"net/http"
	"webSocketsApp1/internal/handlers"
)

func main() {

	// Get the apps routes
	mux := routes()

	log.Println("Starting channel listener 💁")
	go handlers.ListenToWebSocketChannel()

	// Start a web server
	log.Println("Starting a web server on port 8080")
	_ = http.ListenAndServe(":8080", mux)

}
