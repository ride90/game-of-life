package multiverse

import (
	"encoding/json"
	"fmt"
	"github.com/ride90/game-of-life/configs"
	"github.com/ride90/game-of-life/internal/universe"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

type Multiverse struct {
	universes [24]*universe.Universe
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

	// Move all universes to the right in an array (index++).
	for i := len(r.universes) - 1; i >= 0; i-- {
		if r.universes[i] == nil {
			continue
		}
		r.universes[i+1] = r.universes[i]
	}
	r.universes[0] = u
	r.count++
}

func (r *Multiverse) removeStaticUniverses() {

	// // Remove a reference -> garbage collected.
	// r.universes[indexToRemove] = nil
	//

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

func (r *Multiverse) Evolve(cfg *configs.Config) {
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

	// Remove stale static universes.
	indicesToRemove := make([]int, 0, 8)
	for i, u := range r.universes {
		if u != nil && u.IsStatic {
			duration := time.Now().UTC().Sub(u.StaticFrom)
			if cfg.Game.RemoveStaticUniverseAfter <= int(duration.Seconds()) {
				indicesToRemove = append(indicesToRemove, i)
			}
		}
	}
	if len(indicesToRemove) > 0 {
		// Remove references to stale universes -> garbage collected.
		for _, i := range indicesToRemove {
			log.Info("Removing stale static", r.universes[i])
			r.universes[i] = nil
		}
		// Squash left non-nil elements.
		tmpArr := [24]*universe.Universe{}
		tmpIndex := 0
		for i := range r.universes {
			if r.universes[i] == nil {
				continue
			}
			tmpArr[tmpIndex] = r.universes[i]
			tmpIndex++
		}
		r.universes = tmpArr
		r.count = tmpIndex
	}
}

func (r *Multiverse) Reset() {
	log.Infoln("Reset multiverse", r)
	r.universes = [24]*universe.Universe{}
	r.count = 0
}

func (r *Multiverse) ToJSON() ([]byte, error) {
	return json.Marshal(r.universes[:r.count])
}

// Create an empty multiverse.
// This variable will be accessible from multiple places/goroutines.
// Lock is used to avoid a race-conditions.
var mvCreateInstanceLock = &sync.Mutex{}
var mvInstance *Multiverse

func GetInstance() *Multiverse {
	// Singleton anti-pattern is used rather for learning purpose (works fine btw).
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
