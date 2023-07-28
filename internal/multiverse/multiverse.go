package multiverse

import (
	"fmt"
	"github.com/ride90/game-of-life/internal/universe"
	"strings"
	"sync"
)

type Multiverse struct {
	universes []*universe.Universe
}

func newMultiverse() *Multiverse {
	mu := Multiverse{}
	mu.universes = make([]*universe.Universe, 0, 16)
	return &mu
}

func (r *Multiverse) AddUniverse(u *universe.Universe) {
	r.universes = append(r.universes, u)
}

func (r *Multiverse) String() string {
	return fmt.Sprintf("Multiverse with %d universes", len(r.universes))
}

func (r *Multiverse) RenderMatrices() string {
	var matricesStringBuilder strings.Builder
	for i, u := range r.universes {
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
var mvLock = &sync.Mutex{}
var mvInstance *Multiverse

func GetInstance() *Multiverse {
	if mvInstance == nil {
		mvLock.Lock()
		defer mvLock.Unlock()
		mvInstance = newMultiverse()
		if mvInstance == nil {
			mvInstance = newMultiverse()
		}
	}
	return mvInstance
}
