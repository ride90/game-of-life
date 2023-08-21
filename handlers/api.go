package handlers

import (
	"encoding/json"
	"github.com/ride90/game-of-life/configs"
	"github.com/ride90/game-of-life/internal/multiverse"
	"github.com/ride90/game-of-life/internal/universe"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// HandlerAPI API requests handler
type HandlerAPI struct {
	config *configs.Config
}

// NewHandlerAPI creates a new instance of HandlerAPI
func NewHandlerAPI(cfg *configs.Config) HandlerAPI {
	return HandlerAPI{config: cfg}
}

// Health handles the health endpoint request
func (h HandlerAPI) Health(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode("ok")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// CreateUniverse handles the creation of a new universe
func (h HandlerAPI) CreateUniverse(w http.ResponseWriter, r *http.Request) {
	// Get multiverse and ensure we have space for a new universe.
	mv := multiverse.GetInstance()
	if mv.IsFull() {
		log.Warn("Not possible to create universe. Multiverse is full.")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("Multiverse is full")
		return
	}

	// Decode from stream into Universe struct instance.
	var u universe.Universe
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Calculate initial universe stats.
	u.UpdateStats()
	log.Infoln("Created new universe", &u)

	// Add universe into multiverse.
	if h.config.Game.UniversePrepend {
		mv.PrependUniverse(&u)
	} else {
		mv.AppendUniverse(&u)
	}

	// Write response status.
	w.WriteHeader(http.StatusCreated)
}

// ResetMultiverse handles the resetting of the multiverse
func (h HandlerAPI) ResetMultiverse(w http.ResponseWriter, r *http.Request) {
	// Reset multiverse.
	mv := multiverse.GetInstance()
	mv.Reset()

	// Write response status.
	w.WriteHeader(http.StatusOK)
}

// MergeUniverses handles the merging of all universes together
func (h HandlerAPI) MergeUniverses(w http.ResponseWriter, r *http.Request) {

	mv := multiverse.GetInstance()
	mv.Merge()

	// Write response status.
	w.WriteHeader(http.StatusOK)
}
