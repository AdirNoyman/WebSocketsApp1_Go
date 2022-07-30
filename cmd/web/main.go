package main

import (
	"log"
	"net/http"
)

func main() {

	// Get the apps routes
	mux := routes()

	// Start a web server
	log.Println("Starting a web server on port 8080")
	_ = http.ListenAndServe(":8080", mux)

}
