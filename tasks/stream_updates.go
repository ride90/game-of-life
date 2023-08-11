package tasks

import (
	"github.com/ride90/game-of-life/internal/multiverse"
	"github.com/ride90/game-of-life/internal/ws"
	"log"
	"time"
)

func StreamUpdates(wsHub *ws.Hub) {
	// TODO: There should be some kind of lock to see if evolve is busy or not.
	mv := multiverse.GetInstance()
	ticker := time.NewTicker(1000 / 8 * time.Millisecond)
	// ticker := time.NewTicker(2000 * time.Millisecond)

	for _ = range ticker.C {
		// Evolve every universe inside multiverse.
		mv.Evolve()
		// TODO: Stream updates here..
		jsonData, err := mv.ToJSON()
		if err != nil {
			log.Printf("Error while marshaling multiverse into JSON: %s", err)
		}
		wsHub.Broadcast(jsonData)
	}
}
