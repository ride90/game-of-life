package handlers

import (
	"github.com/gorilla/websocket"
	"github.com/ride90/game-of-life/configs"
	"github.com/ride90/game-of-life/internal/ws"
	"log"
	"net/http"
	"time"
)

// HandlerWS handles WebSocket connections
type HandlerWS struct {
	upgrader websocket.Upgrader // Upgrader for upgrading HTTP connections to WebSocket connections
	config   *configs.Config
}

// NewHandlerWS creates a new instance of HandlerWS with the provided configuration
func NewHandlerWS(cfg *configs.Config) HandlerWS {
	return HandlerWS{
		config: cfg,
		upgrader: websocket.Upgrader{
			WriteBufferSize:  cfg.Server.WsWriteBufferSize,
			ReadBufferSize:   cfg.Server.WsReadBufferSize,
			HandshakeTimeout: time.Duration(cfg.Server.WsHandshakeTimeout) * time.Second,
			CheckOrigin:      func(r *http.Request) bool { return true }, // Allow all origins
		},
	}
}

// NewConnection upgrades an HTTP connection to a WebSocket connection and adds it to the hub
func (h HandlerWS) NewConnection(w http.ResponseWriter, r *http.Request, wsHub *ws.Hub) {
	// Upgrade this connection to a WebSocket connection.
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// Add connection to the web socket hub.
	wsConn := &ws.Connection{Conn: conn, Hub: wsHub}
	wsHub.AddConnection(wsConn)

	// Start reading messages from the connection.
	wsConn.ReadMessages()
}
