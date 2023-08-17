package universe

import (
	"fmt"
	"hash/fnv"
	"strings"
	"time"
)

const (
	deadRender  = "x"
	aliveRender = "■"
	deadValue   = false
	aliveValue  = true
)

type Universe struct {
	// TODO: Think of a decomposition json-specific fields.
	//  - https://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
	Matrix           [30][30]bool `json:"cells"`
	Colour           string       `json:"colour"`
	IsStatic         bool         `json:"-"`
	StaticFrom       time.Time    `json:"-"`
	generationNumber int          `json:"-"`
	aliveCellsCount  int          `json:"-"`
	matrixHash       uint64       `json:"-"`
}

func (r *Universe) String() string {
	return fmt.Sprintf(
		"Colour: %s Static: %t Generation %d Alive: %d",
		r.Colour, r.IsStatic, r.generationNumber, r.aliveCellsCount,
	)
}

func (r *Universe) UpdateStats() {
	for _, row := range r.Matrix {
		for _, cell := range row {
			if cell == aliveValue {
				r.aliveCellsCount++
			}
		}
	}
}

func (r *Universe) RenderMatrix() string {
	var matrixStringBuilder strings.Builder
	for _, row := range r.Matrix {
		for _, cell := range row {
			if cell == aliveValue {
				matrixStringBuilder.WriteString(aliveRender)
			} else {
				matrixStringBuilder.WriteString(deadRender)
			}
		}
		matrixStringBuilder.WriteString("\n")
	}
	return matrixStringBuilder.String()
}

func (r *Universe) Evolve() {
	// No sense to compute static universe.
	if r.IsStatic {
		r.generationNumber++
		return
	}

	// Calculate matrix hash & compare with a previous one.
	// If hashes are equal then it means universe is static, and it's not going
	// to evolve, in this case no sense to compute it anymore.
	matrixHash := getMatrixHash(r.Matrix)
	if matrixHash == r.matrixHash {
		r.IsStatic = true
		r.StaticFrom = time.Now().UTC()
		r.generationNumber++
		return
	}
	r.matrixHash = matrixHash

	// Run game of live algorithm.
	nextGenMatrix := r.Matrix
	r.aliveCellsCount = 0
	for y := range r.Matrix {
		for x := range r.Matrix[y] {
			neighborsCount := r.neighboursCount(x, y)
			if r.Matrix[y][x] == aliveValue {
				if neighborsCount < 2 || neighborsCount > 3 {
					nextGenMatrix[y][x] = deadValue
				}
			} else if neighborsCount == 3 {
				nextGenMatrix[y][x] = aliveValue
			}
			if nextGenMatrix[y][x] == aliveValue {
				r.aliveCellsCount++
			}
		}
	}
	r.Matrix = nextGenMatrix
	r.generationNumber++
}

// neighboursCount method calculates the number of live neighbors for a given cell.
func (r *Universe) neighboursCount(x, y int) int {
	var ex, wx, ny, sy, neighbours int
	if x != len(r.Matrix[y])-1 {
		ex = x + 1
	} else {
		ex = 0
	}
	if x != 0 {
		wx = x - 1
	} else {
		wx = len(r.Matrix[y]) - 1
	}
	if y != 0 {
		ny = y - 1
	} else {
		ny = len(r.Matrix) - 1
	}
	if y != len(r.Matrix)-1 {
		sy = y + 1
	} else {
		sy = 0
	}
	if r.Matrix[ny][wx] == aliveValue {
		neighbours++
	}
	if r.Matrix[ny][x] == aliveValue {
		neighbours++
	}
	if r.Matrix[ny][ex] == aliveValue {
		neighbours++
	}
	if r.Matrix[y][wx] == aliveValue {
		neighbours++
	}
	if r.Matrix[y][ex] == aliveValue {
		neighbours++
	}
	if r.Matrix[sy][wx] == aliveValue {
		neighbours++
	}
	if r.Matrix[sy][ex] == aliveValue {
		neighbours++
	}
	if r.Matrix[sy][x] == aliveValue {
		neighbours++
	}
	return neighbours
}

func getMatrixHash(matrix [30][30]bool) uint64 {
	hasher := fnv.New64a()
	for y := range matrix {
		for x := range matrix[y] {
			if matrix[y][x] {
				hasher.Write([]byte{1})
			} else {
				hasher.Write([]byte{0})
			}
		}
	}
	return hasher.Sum64()
}
