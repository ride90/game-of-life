package main

import (
	"fmt"
	"github.com/gorilla/mux"
	config "github.com/ride90/game-of-life"
	"github.com/ride90/game-of-life/handlers"
	"github.com/ride90/game-of-life/internal/ws"
	"github.com/ride90/game-of-life/middlewares"
	"github.com/ride90/game-of-life/tasks"
	"log"

	"net/http"
	"time"
)

// TODO: Create configs https://dev.to/ilyakaznacheev/a-clean-way-to-pass-configs-in-a-go-application-1g64

func main() {
	// Load config.
	// TODO: Think of a better approach of how to auto include config in handlers/tasks.
	cfg := config.NewConfig()

	// Container for websocket connection.
	wsHub := ws.NewHub()

	// Evolve universes & stream updates via ws to clients.
	go tasks.StreamUpdates(wsHub, cfg)

	router := mux.NewRouter()
	// Global middlewares.
	router.Use(middlewares.MiddlewareLogging)

	// API middlewares.
	routerAPI := router.PathPrefix("/api").Subrouter()
	routerAPI.Use(middlewares.MiddlewareContentType)

	// API handlers.
	apiHandler := handlers.NewHandlerAPI(cfg)
	routerAPI.HandleFunc("/health", apiHandler.Health).Methods(http.MethodGet)
	routerAPI.HandleFunc("/universe", apiHandler.CreateUniverse).Methods(http.MethodPost)

	// WS handler.
	wsHandler := handlers.NewHandlerWS(cfg)
	router.HandleFunc(
		"/ws/updates",
		func(w http.ResponseWriter, r *http.Request) {
			wsHandler.NewConnection(w, r, wsHub)
		},
	)

	// Static files handler.
	spaHandler := handlers.NewHandlerSPA("web", "index.html")
	router.PathPrefix("/").Handler(spaHandler)
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	// Create & run server.
	srv := &http.Server{
		Handler: router,
		Addr:    addr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
	}
	log.Println("Running", addr)
	log.Fatal(srv.ListenAndServe())
}
