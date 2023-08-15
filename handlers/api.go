package handlers

import (
	"encoding/json"
	config "github.com/ride90/game-of-life"
	"github.com/ride90/game-of-life/internal/multiverse"
	"github.com/ride90/game-of-life/internal/universe"
	"net/http"
)

type HandlerAPI struct {
	config *config.Config
}

func NewHandlerAPI(cfg *config.Config) HandlerAPI {
	return HandlerAPI{config: cfg}
}

func (h HandlerAPI) Health(w http.ResponseWriter, r *http.Request) {
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
	if h.config.Game.UniversePrepend {
		mv.PrependUniverse(&u)
	} else {
		mv.AppendUniverse(&u)
	}

	// Write response status.
	w.WriteHeader(http.StatusCreated)
}
