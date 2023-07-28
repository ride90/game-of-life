package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ride90/game-of-life/internal/multiverse"
	"github.com/ride90/game-of-life/internal/universe"
	"net/http"
)

type handlerAPI struct{}

func NewHandlerAPI() handlerAPI {
	return handlerAPI{}
}

func (h handlerAPI) Health(w http.ResponseWriter, r *http.Request) {

	mv := multiverse.GetInstance()
	fmt.Println(mv.RenderMatrices())

	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func (h handlerAPI) CreateUniverse(w http.ResponseWriter, r *http.Request) {
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
	mv := multiverse.GetInstance()
	mv.AddUniverse(&u)

	w.WriteHeader(http.StatusCreated)
}
