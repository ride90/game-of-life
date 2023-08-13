package multiverse

import (
	"encoding/json"
	"fmt"
	"github.com/ride90/game-of-life/internal/universe"
	"strings"
	"sync"
)

type Multiverse struct {
	universes [32]*universe.Universe
	count     int
	// This lock will be used during evolve + AppendUniverse.
	// In theory, during evolve this lock is not needed, since every universe
	// has its own memory address and concurrent evolving of universes is totally
	// race-condition free. Even AppendUniverse should be safe since it's only
	// "appending" a universe to an array and there is not much of intersection
	// with evolve, but it's better to cover AppendUniverse with a lock.
	lock sync.Mutex
}

func newMultiverse() *Multiverse {
	mu := Multiverse{}
	return &mu
}

func (r *Multiverse) AppendUniverse(u *universe.Universe) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.universes[r.count] = u
	r.count++
}

func (r *Multiverse) PrependUniverse(u *universe.Universe) {
	r.lock.Lock()
	defer r.lock.Unlock()

	// Ensure we can fit a new universe.
	if r.IsFull() {
		return
	}

	// Move all universes to the left in an array (index++).
	for i := len(r.universes) - 1; i >= 0; i-- {
		if r.universes[i] == nil {
			continue
		}
		r.universes[i+1] = r.universes[i]
	}
	r.universes[0] = u
	r.count++
}

func (r *Multiverse) IsFull() bool {
	return r.count >= len(r.universes)
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

func (r *Multiverse) Evolve() {
	// Lock & Unlock.
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()
	// Each universe evolves itself in goroutine.
	var wg sync.WaitGroup
	for _, u := range r.universes {
		if u == nil {
			continue
		}
		// TODO: Think/research if closure approach is better here:
		//   https://go.dev/doc/faq#closures_and_goroutines
		wg.Add(1)
		go func(u *universe.Universe, wg *sync.WaitGroup) {
			defer wg.Done()
			u.Evolve()
		}(u, &wg)
	}
	wg.Wait()
}

func (r *Multiverse) ToJSON() ([]byte, error) {
	return json.Marshal(r.universes[:r.count])
}

// Create an empty multiverse.
// This variable will be accessible from multiple places/goroutines.
// Lock is used to avoid a race-conditions.
// Singleton anti-pattern is used rather for learning purpose (works fine btw).
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
