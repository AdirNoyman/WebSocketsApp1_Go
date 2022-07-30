package main

import (
	"github.com/bmizerany/pat"
	"net/http"
	"webSocketsApp1/internal/handlers"
)

func routes() http.Handler {

	// Create an instance of the PAT route handler
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Home))
	mux.Get("/ws", http.HandlerFunc(handlers.WebSocketEndPoint))

	return mux

}
