package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type handlerAPI struct{}

func NewHandlerAPI() handlerAPI {
	return handlerAPI{}
}

func (h handlerAPI) Health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func (h handlerAPI) CreateUniverse(w http.ResponseWriter, r *http.Request) {

	fmt.Println("!!!http.Request", r)

	w.WriteHeader(http.StatusCreated)
}
