package tasks

import (
	"github.com/ride90/game-of-life/internal/multiverse"
	"github.com/ride90/game-of-life/internal/ws"
	"log"
	"time"
)

func StreamUpdates(wsHub *ws.Hub) {
	mv := multiverse.GetInstance()
	ticker := time.NewTicker(1000 / 5 * time.Millisecond)
	locked := false

	for _ = range ticker.C {
		// Check if operation is locked.
		if locked {
			continue
		}
		locked = true

		// Evolve every universe inside multiverse.
		mv.Evolve()
		// Prepare json and broadcast it to all ws client.
		jsonData, err := mv.ToJSON()
		if err != nil {
			log.Printf("Error while marshaling multiverse into JSON: %s", err)
		}
		wsHub.Broadcast(jsonData)

		// Unlock.
		locked = false
	}
}
