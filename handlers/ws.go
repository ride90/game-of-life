package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/ride90/game-of-life/internal/ws"
	"log"
	"net/http"
)

type HandlerWS struct {
	upgrader websocket.Upgrader
}

func NewHandlerWS() HandlerWS {
	return HandlerWS{
		upgrader: websocket.Upgrader{
			// TODO: Set a reasonable buffer.
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (h HandlerWS) NewConnection(w http.ResponseWriter, r *http.Request, wsHub *ws.Hub) {
	// Upgrade this connection to a WebSocket connection.
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// Add connection to the web socket hub.
	wsConn := &ws.Connection{Conn: conn, Hub: wsHub}
	wsHub.AddConnection(wsConn)

	// Start reading messages.
	wsConn.ReadMessages()
}
