package logger

import (
	"github.com/ride90/game-of-life/configs"
	log "github.com/sirupsen/logrus"
	"os"
)

func SetupLogger(cfg *configs.Config) {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:            true,
		DisableLevelTruncation: true,
		FullTimestamp:          true,
	})
	log.SetOutput(os.Stdout)
	log.SetReportCaller(cfg.Log.SetReportCaller)

	// Get log level from the config.
	level, err := log.ParseLevel(cfg.Log.Level)
	if err != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
}
