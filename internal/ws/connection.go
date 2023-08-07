package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Connection struct {
	Conn             *websocket.Conn
	Hub              *Hub
	readMessagesLock sync.Mutex
}

func (r *Connection) String() string {
	return fmt.Sprintf(
		"WS Connection. Remote: %s. Local: %s",
		r.Conn.RemoteAddr(), r.Conn.LocalAddr(),
	)
}

func (r *Connection) ReadMessages() {
	r.readMessagesLock.Lock()
	defer r.readMessagesLock.Unlock()

	// Read messages.
	for {
		_, _, err := r.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Error reading message: %v", err)
			}
			log.Println("Connection closed by the client.")
			break
		}
	}
	// Connection closed.
	r.Hub.RemoveConnection(r)
}
