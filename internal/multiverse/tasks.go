package multiverse

import (
	"fmt"
	"time"
)

func EvolveMultiverseTask() {
	mv := GetInstance()
	ticker := time.NewTicker(5 * time.Second)
	for _ = range ticker.C {
		fmt.Println("Evolving", mv)
	}
}
