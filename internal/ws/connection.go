package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"sync"
)

// Connection represents a WebSocket connection.
type Connection struct {
	Conn             *websocket.Conn // The underlying WebSocket connection
	Hub              *Hub
	readMessagesLock sync.Mutex
}

// String returns a formatted string representation of the connection.
func (r *Connection) String() string {
	return fmt.Sprintf(
		"WS Connection. Remote: %s. Local: %s",
		r.Conn.RemoteAddr(), r.Conn.LocalAddr(),
	)
}

// ReadMessages reads messages from the WebSocket connection.
func (r *Connection) ReadMessages() {
	// Lock to ensure single-threaded message reading
	r.readMessagesLock.Lock()
	defer r.readMessagesLock.Unlock()

	// Read messages.
	for {
		_, _, err := r.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Error("Error reading message: %v", err)
			}
			log.Debug("Connection closed by the client.")
			break
		}
	}

	// Connection closed, remove from the hub.
	r.Hub.RemoveConnection(r)
}

// SendMessage sends a WebSocket message with the provided data.
func (r *Connection) SendMessage(data []byte) {
	err := r.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Error("Error while sending a message. %s. Error: %s", r, err)
	}
}
