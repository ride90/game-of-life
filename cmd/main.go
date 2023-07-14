package main

import (
	"github.com/gorilla/mux"
	"github.com/ride90/game-of-life/handlers"
	"log"
	"net/http"
)

func main() {
	handler := handlers.New()
	router := mux.NewRouter()
	router.HandleFunc("/", handler.Root).Methods(http.MethodGet)

	log.Println("Game of Life is running!")
	http.ListenAndServe(":4000", router)
}
