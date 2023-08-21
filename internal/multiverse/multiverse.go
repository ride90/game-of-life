package multiverse

import (
	"encoding/json"
	"fmt"
	"github.com/ride90/game-of-life/configs"
	"github.com/ride90/game-of-life/internal/universe"
	log "github.com/sirupsen/logrus"
	"math"
	"strings"
	"sync"
	"time"
)

const (
	universesPerRow      = 4
	mergedUniverseColour = "#F00"
)

var mvCreateInstanceLock = &sync.Mutex{}
var mvInstance *Multiverse

// GetInstance retrieves or creates a singleton instance of Multiverse
// Singleton anti-pattern is used here for learning purposes.
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

// Multiverse represents the collection of universes
type Multiverse struct {
	universes [24]*universe.Universe
	count     int
	lock      sync.Mutex // Mutex for concurrent access control
}

// newMultiverse creates a new instance of Multiverse
func newMultiverse() *Multiverse {
	mu := Multiverse{}
	return &mu
}

// AppendUniverse adds a new universe to the end of the collection
func (r *Multiverse) AppendUniverse(u *universe.Universe) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.universes[r.count] = u
	r.count++
}

// PrependUniverse adds a new universe to the beginning of the collection
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

// IsFull checks if the Multiverse is full
func (r *Multiverse) IsFull() bool {
	return r.count >= len(r.universes)
}

// String returns a string representation of the Multiverse
func (r *Multiverse) String() string {
	return fmt.Sprintf("Multiverse with %d/%d universes", r.count, len(r.universes))
}

// RenderMatrices returns a string containing rendered matrices of contained universes
// Used for stdout & debug purposes.
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

// Evolve evolves all universes in the Multiverse
func (r *Multiverse) Evolve(cfg *configs.Config) {
	// Lock & Unlock.
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()

	// Each universe evolves itself in a goroutine.
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

// Reset clears the Multiverse
func (r *Multiverse) Reset() {
	log.Infoln("Reset multiverse", r)
	r.universes = [24]*universe.Universe{}
	r.count = 0
}

// Merge merges all universes together into one big madness
func (r *Multiverse) Merge() {
	// Lock & Unlock.
	r.lock.Lock()
	defer func() {
		r.lock.Unlock()
	}()

	// Check if it makes sense to perform merge
	if r.count <= 1 {
		log.Warn("Merge doesn't make sense", r)
		return
	}

	log.Infoln("Performing universes merge", r)

	// Create matrix to fit all universes.
	normalMatrixSize := len(r.universes[0].Matrix[0])
	var finalMatrix [][]bool
	finalMatrixHeight := int(math.Ceil(float64(r.count)/universesPerRow)) * normalMatrixSize
	finalMatrixWidth := normalMatrixSize * universesPerRow
	finalMatrix = make([][]bool, finalMatrixHeight)
	for i := range finalMatrix {
		finalMatrix[i] = make([]bool, finalMatrixWidth)
	}

	// Fit "normal" matrices into a "big" final one.
	// Having 3 nested loops is fine here since we are merging array of
	// into one matrix. It's happening only once when merge is triggered.
	var finalX, finalY = 0, 0
	for indexUniverse, universe := range r.universes {
		if universe == nil {
			break
		}
		if indexUniverse%universesPerRow == 0 {
			finalX = 0
		} else {
			finalY = indexUniverse / universesPerRow * normalMatrixSize
		}

		for y := range universe.Matrix {
			for x := range universe.Matrix[y] {
				finalMatrix[finalY][finalX] = universe.Matrix[y][x]
				finalX++
			}
			finalY++
			finalX = (indexUniverse % universesPerRow) * 50
		}
	}

	// Create a universe which will keep all universes inside
	finalUniverse := universe.Universe{
		Colour: mergedUniverseColour,
		Matrix: finalMatrix,
	}

	// Reset multiverse -> remove all universes & reset count
	r.Reset()

	// Add final universe
	r.universes[r.count] = &finalUniverse
	r.count++
}

// ToJSON serializes the Multiverse to JSON format
func (r *Multiverse) ToJSON() ([]byte, error) {
	return json.Marshal(r.universes[:r.count])
}
