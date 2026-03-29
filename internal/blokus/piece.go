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

func getPiecesAtStart() []*piece {
	pieces := []*piece{
		// 1 Monomino
		{id: 0, cells: []coordinate{{0, 0}}},
		// 1 Domino
		{id: 1, cells: []coordinate{{0, 0}, {1, 0}}},
		// 2 Trominoes
		{id: 2, cells: []coordinate{{0, 0}, {1, 0}, {2, 0}}},
		{id: 3, cells: []coordinate{{0, 0}, {1, 0}, {0, 1}}},
		// 5 Tetrominoes
		{id: 4, cells: []coordinate{{0, 0}, {1, 0}, {0, 1}, {1, 1}}},
		{id: 5, cells: []coordinate{{0, 0}, {1, 0}, {2, 0}, {1, 1}}},
		{id: 6, cells: []coordinate{{0, 0}, {1, 0}, {2, 0}, {0, 1}}},
		{id: 7, cells: []coordinate{{0, 0}, {1, 0}, {1, 1}, {2, 1}}},
		{id: 8, cells: []coordinate{{0, 0}, {1, 0}, {2, 0}, {3, 0}}},
		// 12 Pentominoes
		{id: 9, cells: []coordinate{{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}}},
		{id: 10, cells: []coordinate{{0, 0}, {1, 0}, {2, 0}, {3, 0}, {0, 1}}},
		{id: 11, cells: []coordinate{{0, 0}, {1, 0}, {2, 0}, {3, 0}, {1, 1}}},
		{id: 12, cells: []coordinate{{0, 0}, {1, 0}, {1, 1}, {2, 1}, {3, 1}}},
		{id: 13, cells: []coordinate{{0, 0}, {1, 0}, {0, 1}, {1, 1}, {0, 2}}},
		{id: 14, cells: []coordinate{{0, 0}, {2, 0}, {0, 1}, {1, 1}, {2, 1}}},
		{id: 15, cells: []coordinate{{0, 0}, {1, 0}, {2, 0}, {0, 1}, {0, 2}}},
		{id: 16, cells: []coordinate{{0, 0}, {0, 1}, {1, 1}, {1, 2}, {2, 2}}},
		{id: 17, cells: []coordinate{{0, 0}, {1, 0}, {1, 1}, {1, 2}, {2, 2}}},
		{id: 18, cells: []coordinate{{1, 0}, {2, 0}, {0, 1}, {1, 1}, {1, 2}}},
		{id: 19, cells: []coordinate{{1, 0}, {0, 1}, {1, 1}, {2, 1}, {1, 2}}},
		{id: 20, cells: []coordinate{{0, 0}, {1, 0}, {2, 0}, {1, 1}, {1, 2}}},
	}
	for _, p := range pieces {
		p.normalize()
	}
	return pieces
}
