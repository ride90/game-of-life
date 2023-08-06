package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ride90/game-of-life/handlers"
	"github.com/ride90/game-of-life/internal/multiverse"
	"github.com/ride90/game-of-life/middlewares"
	"log"

	"net/http"
	"time"
)

func main() {
	// Start evolving multiverse in the background.
	go multiverse.EvolveMultiverseTask()

	router := mux.NewRouter()
	// Global middlewares.
	router.Use(middlewares.MiddlewareLogging)

	// API middlewares.
	routerAPI := router.PathPrefix("/api").Subrouter()
	routerAPI.Use(middlewares.MiddlewareContentType)

	// API handlers.
	apiHandler := handlers.NewHandlerAPI()
	routerAPI.HandleFunc("/health", apiHandler.Health).Methods(http.MethodGet)
	routerAPI.HandleFunc("/universe", apiHandler.CreateUniverse).Methods(http.MethodPost)

	// WS handler.
	routerWS := router.PathPrefix("/ws").Subrouter()
	wsHandler := handlers.NewHandlerWS()
	routerWS.HandleFunc("/updates", wsHandler.StreamUpdates)

	// Static files handler.
	spaHandler := handlers.NewHandlerSPA("web", "index.html")
	router.PathPrefix("/").Handler(spaHandler)
	// TODO: Move this crap somewhere.
	const host, port = "127.0.0.1", 4000
	addr := fmt.Sprintf("%s:%d", host, port)

	// Create & run server.
	srv := &http.Server{
		Handler: router,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Running", addr)
	log.Fatal(srv.ListenAndServe())
}
