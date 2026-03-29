package blokus

import (
	"math"
	"slices"
)

type Coordinate struct {
	X, Y int
}

type Piece struct {
	Id    int
	Cells []Coordinate
}

// Rotate 90 degree clockwise
// Need to call Normalize afterwards
func (p *Piece) Rotate() {
	for i := range p.Cells {
		c := &p.Cells[i]
		c.X, c.Y = -c.Y, c.X
	}
}

// Flip along the Y axis
// Need to call Normalize afterwards
func (p *Piece) Flip() {
	for i := range p.Cells {
		c := &p.Cells[i]
		c.X = -c.X
	}
}

// Normalize first makes sure all coordinates are non negative
// Then it sorts by X first and by Y to compare to coordinates
func (p *Piece) Normalize() {
	minX, minY := math.MaxInt32, math.MaxInt32
	for _, c := range p.Cells {
		minX = min(minX, c.X)
		minY = min(minY, c.Y)
	}
	for i := range p.Cells {
		c := &p.Cells[i]
		c.X -= minX
		c.Y -= minY
	}
	slices.SortFunc(p.Cells, func(a, b Coordinate) int {
		if a.X == b.X {
			return a.Y - b.Y
		}
		return a.X - b.X
	})
}
