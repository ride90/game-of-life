package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type HandlerWS struct {
	upgrader websocket.Upgrader
}

func NewHandlerWS() HandlerWS {
	return HandlerWS{
		upgrader: websocket.Upgrader{
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (h HandlerWS) StreamUpdates(w http.ResponseWriter, r *http.Request) {
	// Upgrade this connection to a WebSocket connection.
	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(ws)
}
