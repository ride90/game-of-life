package universe

import (
	"fmt"
	"strings"
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
	generationNumber int          `json:"-"`
	aliveCellsCount  int          `json:"-"`
}

func (r *Universe) String() string {
	return fmt.Sprintf(
		"Colour: %s\nGeneration %d\nAlive: %d\n%s",
		r.Colour, r.generationNumber, r.aliveCellsCount, r.RenderMatrix(),
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
