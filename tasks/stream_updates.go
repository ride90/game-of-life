package tasks

import (
	"fmt"
	"github.com/ride90/game-of-life/internal/multiverse"
	"github.com/ride90/game-of-life/internal/ws"
	"time"
)

func StreamUpdates(wsHub *ws.Hub) {
	// TODO: There should be some kind of lock to see if evolve is busy or not.
	mv := multiverse.GetInstance()
	// ticker := time.NewTicker(1000 / 25 * time.Millisecond)
	ticker := time.NewTicker(10000 * time.Millisecond)

	for _ = range ticker.C {
		// Evolve every universe inside multiverse.
		mv.Evolve()
		// TODO: Stream updates here..
		fmt.Println(wsHub)
	}
}
