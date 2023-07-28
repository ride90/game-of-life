package multiverse

import (
	"fmt"
	"github.com/ride90/game-of-life/internal/universe"
	"strings"
	"sync"
)

type Multiverse struct {
	universes [16]*universe.Universe
	count     int
	// This lock will be used during evolve + add new universe.
	// In theory during evolve this lock is not needed.
	// TODO: Implement lock.
	writeLock sync.Mutex
}

func newMultiverse() *Multiverse {
	mu := Multiverse{}
	return &mu
}

func (r *Multiverse) AddUniverse(u *universe.Universe) {
	r.writeLock.Lock()
	defer r.writeLock.Unlock()
	r.universes[r.count] = u
	r.count++
}

func (r *Multiverse) String() string {
	return fmt.Sprintf("Multiverse with %d/%d universes", r.count, len(r.universes))
}

func (r *Multiverse) RenderMatrices() string {
	var matricesStringBuilder strings.Builder
	for i, u := range r.universes {
		if u == nil {
			continue
		}
		matricesStringBuilder.WriteString(
			fmt.Sprintf("Matrix #%d:\n", i),
		)
		matricesStringBuilder.WriteString(
			fmt.Sprintf("%s\n", u.RenderMatrix()),
		)
	}
	return matricesStringBuilder.String()
}

// Create an empty multiverse.
// This variable will be accessible from multiple places/goroutines.
// Lock is used to avoid a race-conditions.
// Singleton anti-pattern is used.
var mvCreateInstanceLock = &sync.Mutex{}
var mvInstance *Multiverse

func GetInstance() *Multiverse {
	if mvInstance == nil {
		mvCreateInstanceLock.Lock()
		defer mvCreateInstanceLock.Unlock()
		mvInstance = newMultiverse()
		if mvInstance == nil {
			mvInstance = newMultiverse()
		}
	}
	return mvInstance
}
