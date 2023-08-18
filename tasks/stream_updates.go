package tasks

import (
	"github.com/ride90/game-of-life/configs"
	"github.com/ride90/game-of-life/internal/multiverse"
	"github.com/ride90/game-of-life/internal/ws"
	log "github.com/sirupsen/logrus"
	"time"
)

func StreamUpdates(wsHub *ws.Hub, cfg *configs.Config) {
	mv := multiverse.GetInstance()
	ticker := time.NewTicker(1000 / time.Duration(cfg.Game.Fps) * time.Millisecond)
	locked := false

	for _ = range ticker.C {
		// Check if operation is locked.
		if locked {
			continue
		}
		locked = true

		// Evolve every universe inside multiverse.
		mv.Evolve(cfg)
		// Prepare json and broadcast it to all ws client.
		jsonData, err := mv.ToJSON()
		if err != nil {
			log.Errorf("Error while marshaling multiverse into JSON: %s", err)
		}
		wsHub.Broadcast(jsonData)

		// Unlock.
		locked = false
	}
}
