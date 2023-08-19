package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// Hub represents a WebSocket hub that manages connections.
type Hub struct {
	connections []*Connection
}

// NewHub creates a new instance of a WebSocket hub.
func NewHub() *Hub {
	return &Hub{
		connections: make([]*Connection, 0, 16),
	}
}

// String returns a formatted string representation of the hub.
func (r *Hub) String() string {
	return fmt.Sprintf(
		"WS Hub. Active connections: %d", len(r.connections),
	)
}

// AddConnection adds a new WebSocket connection to the hub.
func (r *Hub) AddConnection(c *Connection) {
	log.Debug("Adding new", c)
	r.connections = append(r.connections, c)
}

// RemoveConnection removes a WebSocket connection from the hub.
func (r *Hub) RemoveConnection(c *Connection) bool {
	for i := 0; i < len(r.connections); i++ {
		if r.connections[i] == c {
			// Try to close the connection on our side.
			defer func(Conn *websocket.Conn) {
				err := Conn.Close()
				if err != nil {
					log.Debug("Closing WS connection:", err)
				}
			}(r.connections[i].Conn)

			// Remove from the hub.
			log.Debug("WS Hub removing:", r.connections[i])
			r.connections[i] = r.connections[len(r.connections)-1]
			r.connections = r.connections[:len(r.connections)-1]
			return true
		}
	}
	return false
}

// Broadcast sends a message to all clients connected to the hub.
func (r *Hub) Broadcast(data []byte) {
	log.Debugf("Broadcasting message to %d clients", len(r.connections))
	for _, connection := range r.connections {
		connection.SendMessage(data)
	}
}
