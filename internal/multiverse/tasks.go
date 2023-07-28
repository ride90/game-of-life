package multiverse

import (
	"fmt"
	"time"
)

func EvolveMultiverseTask() {
	mv := GetInstance()
	ticker := time.NewTicker(5 * time.Second)
	for _ = range ticker.C {
		mv.writeLock.Lock()
		fmt.Println("Evolving", mv)
		time.Sleep(2 * time.Second)
		mv.writeLock.Unlock()
	}
}
