package multiverse

import (
	"time"
)

func EvolveMultiverseTask() {
	// TODO: There should be some kind of lock to see if evolve is busy or not.
	mv := GetInstance()
	ticker := time.NewTicker(300 * time.Millisecond)
	for _ = range ticker.C {
		mv.evolve()
	}
}
