package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Hub struct {
	connections []*Connection
}

func NewHub() *Hub {
	return &Hub{
		connections: make([]*Connection, 0, 16),
	}
}

func (r *Hub) String() string {
	return fmt.Sprintf(
		"WS Hub. Active connections: %d", len(r.connections),
	)
}

func (r *Hub) AddConnection(c *Connection) {
	log.Debug("Adding new", c)
	r.connections = append(r.connections, c)
}

func (r *Hub) RemoveConnection(c *Connection) bool {
	for i := 0; i < len(r.connections); i++ {
		if r.connections[i] == c {
			// Try to close connection on our side.
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

func (r *Hub) Broadcast(data []byte) {
	log.Debugf("Broadcasting message to %d clients", len(r.connections))
	for _, connection := range r.connections {
		connection.SendMessage(data)
	}
}
