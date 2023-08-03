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
}
