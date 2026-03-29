package blokus

import (
	"math"
	"slices"
)

type coordinate struct {
	x, y int
}

type piece struct {
	id    int
	cells []coordinate
}

// Rotate 90 degree clockwise
// Need to call Normalize afterwards
func (p *piece) rotate() {
	for i := range p.cells {
		c := &p.cells[i]
		c.x, c.y = -c.y, c.x
	}
}

// Flip along the Y axis
// Need to call Normalize afterwards
func (p *piece) flip() {
	for i := range p.cells {
		c := &p.cells[i]
		c.x = -c.x
	}
}

// Normalize first makes sure all coordinates are non negative
// Then it sorts by X first and by Y to compare to coordinates
func (p *piece) normalize() {
	minX, minY := math.MaxInt32, math.MaxInt32
	for _, c := range p.cells {
		minX = min(minX, c.x)
		minY = min(minY, c.y)
	}
	for i := range p.cells {
		c := &p.cells[i]
		c.x -= minX
		c.y -= minY
	}
	slices.SortFunc(p.cells, func(a, b coordinate) int {
		if a.x == b.x {
			return a.y - b.y
		}
		return a.x - b.x
	})
}
