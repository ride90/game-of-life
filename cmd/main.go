package main

import (
	"fmt"
	"github.com/ride90/game-of-life/handlers"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// API handlers.
	apiHandler := handlers.NewHandlerAPI()
	router.HandleFunc("/api/health", apiHandler.Health).Methods(http.MethodGet)

	// Static files handler.
	spaHandler := handlers.NewHandlerSPA("static", "index.html")
	router.PathPrefix("/").Handler(spaHandler)

	const host, port = "127.0.0.1", 4000
	addr := fmt.Sprintf("%s:%d", host, port)

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
