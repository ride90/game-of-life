package handlers

import (
	"github.com/gorilla/websocket"
	config "github.com/ride90/game-of-life"
	"github.com/ride90/game-of-life/internal/ws"
	"log"
	"net/http"
	"time"
)

type HandlerWS struct {
	upgrader websocket.Upgrader
	config   *config.Config
}

func NewHandlerWS(cfg *config.Config) HandlerWS {
	return HandlerWS{
		config: cfg,
		upgrader: websocket.Upgrader{
			WriteBufferSize:  cfg.Server.WsWriteBufferSize,
			ReadBufferSize:   cfg.Server.WsReadBufferSize,
			HandshakeTimeout: time.Duration(cfg.Server.WsHandshakeTimeout) * time.Second,
			CheckOrigin:      func(r *http.Request) bool { return true },
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
