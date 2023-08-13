package handlers

import (
	"encoding/json"
	"github.com/ride90/game-of-life/internal/multiverse"
	"github.com/ride90/game-of-life/internal/universe"
	"net/http"
)

type HandlerAPI struct{}

func NewHandlerAPI() HandlerAPI {
	return HandlerAPI{}
}

func (h HandlerAPI) Health(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement a proper health check.
	err := json.NewEncoder(w).Encode("ok")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h HandlerAPI) CreateUniverse(w http.ResponseWriter, r *http.Request) {
	// Get multiverse and ensure we have a space for a new universe.
	mv := multiverse.GetInstance()
	if mv.IsFull() {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Multiverse is full")
		return
	}

	var u universe.Universe
	// Decode from stream into Universe struct instance.
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Calculate initial universe stats.
	u.UpdateStats()
	// Add universe into multiverse.
	// TODO: Make append & prepend approach configurable.
	mv.PrependUniverse(&u)
	// Write response status.
	w.WriteHeader(http.StatusCreated)
}
