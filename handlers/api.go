package handlers

import (
	"encoding/json"
	"net/http"
)

type handlerAPI struct{}

func NewHandlerAPI() handlerAPI {
	return handlerAPI{}
}

func (h handlerAPI) Health(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
