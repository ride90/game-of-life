package multiverse

import (
	"time"
)

func EvolveMultiverseTask() {
	// TODO: There should be some kind of lock to see if evolve is busy or not.
	mv := GetInstance()
	ticker := time.NewTicker(5 * time.Second)
	for _ = range ticker.C {
		mv.evolve()
	}
}
