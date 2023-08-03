package universe

import (
	"fmt"
	"strings"
)

const (
	dead  = "x"
	alive = "■"
)

type Universe struct {
	Matrix           [20][20]bool `json:"cells"`
	Colour           string       `json:"colour"`
	generationNumber int
	aliveCellsCount  int
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
			if cell {
				r.aliveCellsCount++
			}
		}
	}
}

func (r *Universe) RenderMatrix() string {
	var matrixStringBuilder strings.Builder
	for _, row := range r.Matrix {
		for _, cell := range row {
			if cell {
				matrixStringBuilder.WriteString(alive)
			} else {
				matrixStringBuilder.WriteString(dead)
			}
		}
		matrixStringBuilder.WriteString("\n")
	}
	return matrixStringBuilder.String()
}

func (r *Universe) Evolve() {
	fmt.Println("Evolving universe", r)
	nextGenMatrix := r.Matrix

	r.aliveCellsCount = 0
	for y := range r.Matrix {
		for x := range r.Matrix[y] {
			neighborsCount := r.neighboursCount(x, y)
			if r.Matrix[y][x] {
				if neighborsCount < 2 || neighborsCount > 3 {
					nextGenMatrix[y][x] = false
				}
			} else if neighborsCount == 3 {
				nextGenMatrix[y][x] = true
			}
			if nextGenMatrix[y][x] {
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
	if r.Matrix[ny][wx] {
		neighbours++
	}
	if r.Matrix[ny][x] {
		neighbours++
	}
	if r.Matrix[ny][ex] {
		neighbours++
	}
	if r.Matrix[y][wx] {
		neighbours++
	}
	if r.Matrix[y][ex] {
		neighbours++
	}
	if r.Matrix[sy][wx] {
		neighbours++
	}
	if r.Matrix[sy][ex] {
		neighbours++
	}
	if r.Matrix[sy][x] {
		neighbours++
	}
	return neighbours
}
