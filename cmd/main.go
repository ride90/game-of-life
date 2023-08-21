package main

// TODO: Refactor
// 	- Add return/handle errors in places where it makes sense.

import (
	"fmt"
	"github.com/gorilla/mux"
	logger "github.com/ride90/game-of-life"
	"github.com/ride90/game-of-life/configs"
	"github.com/ride90/game-of-life/handlers"
	"github.com/ride90/game-of-life/internal/ws"
	"github.com/ride90/game-of-life/middlewares"
	"github.com/ride90/game-of-life/tasks"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var cfg *configs.Config
var wsHub *ws.Hub

func init() {
	// Config.
	// TODO: Think of a better approach how to include config in handlers/tasks/internals.
	cfg = configs.NewConfig()

	// Manager for WS connection.
	wsHub = ws.NewHub()

	// Setup a logger.
	logger.SetupLogger(cfg)
}

func main() {
	// Evolve universes & stream updates via ws to clients.
	go tasks.StreamUpdates(wsHub, cfg)

	router := mux.NewRouter()
	// Global middlewares.
	router.Use(middlewares.MiddlewareLogRequest)

	// API middlewares.
	routerAPI := router.PathPrefix("/api").Subrouter()
	routerAPI.Use(middlewares.MiddlewareContentType)

	// API handlers.
	apiHandler := handlers.NewHandlerAPI(cfg)
	routerAPI.HandleFunc("/health", apiHandler.Health).Methods(http.MethodGet)
	routerAPI.HandleFunc("/universe", apiHandler.CreateUniverse).Methods(http.MethodPost)
	routerAPI.HandleFunc("/bigbang", apiHandler.ResetMultiverse).Methods(http.MethodPost)
	routerAPI.HandleFunc("/merge", apiHandler.MergeUniverses).Methods(http.MethodPost)

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
	log.Info("Running server on ", addr)
	log.Fatal(srv.ListenAndServe())
}
