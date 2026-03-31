package blokus

import (
	"math"
	"slices"
)

type Coordinate struct {
	x, y int
}

func NewCoordinate(x, y int) Coordinate {
	return Coordinate{x: x, y: y}
}

func (c Coordinate) X() int {
	return c.x
}

func (c Coordinate) Y() int {
	return c.y
}

type piece struct {
	id    int
	cells []Coordinate
}

// Rotate 90 degree clockwise.
// Need to call Normalize afterwards.
func (p *piece) rotate() {
	for i := range p.cells {
		c := &p.cells[i]
		c.x, c.y = -c.y, c.x
	}
}

// Flip along the Y axis.
// Need to call Normalize afterwards.
func (p *piece) flip() {
	for i := range p.cells {
		c := &p.cells[i]
		c.x = -c.x
	}
}

// Normalize first makes sure all Coordinates are non negative.
// Then it sorts by X first and by Y to compare to Coordinates.
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
	slices.SortFunc(p.cells, func(a, b Coordinate) int {
		if a.x == b.x {
			return a.y - b.y
		}
		return a.x - b.x
	})
}

func getPiecesAtStart() []*piece {
	pieces := []*piece{
		// 1 Monomino
		{id: 0, cells: []Coordinate{{0, 0}}},
		// 1 Domino
		{id: 1, cells: []Coordinate{{0, 0}, {1, 0}}},
		// 2 Trominoes
		{id: 2, cells: []Coordinate{{0, 0}, {1, 0}, {2, 0}}},
		{id: 3, cells: []Coordinate{{0, 0}, {1, 0}, {0, 1}}},
		// 5 Tetrominoes
		{id: 4, cells: []Coordinate{{0, 0}, {1, 0}, {0, 1}, {1, 1}}},
		{id: 5, cells: []Coordinate{{0, 0}, {1, 0}, {2, 0}, {1, 1}}},
		{id: 6, cells: []Coordinate{{0, 0}, {1, 0}, {2, 0}, {0, 1}}},
		{id: 7, cells: []Coordinate{{0, 0}, {1, 0}, {1, 1}, {2, 1}}},
		{id: 8, cells: []Coordinate{{0, 0}, {1, 0}, {2, 0}, {3, 0}}},
		// 12 Pentominoes
		{id: 9, cells: []Coordinate{{0, 0}, {1, 0}, {2, 0}, {3, 0}, {4, 0}}},
		{id: 10, cells: []Coordinate{{0, 0}, {1, 0}, {2, 0}, {3, 0}, {0, 1}}},
		{id: 11, cells: []Coordinate{{0, 0}, {1, 0}, {2, 0}, {3, 0}, {1, 1}}},
		{id: 12, cells: []Coordinate{{0, 0}, {1, 0}, {1, 1}, {2, 1}, {3, 1}}},
		{id: 13, cells: []Coordinate{{0, 0}, {1, 0}, {0, 1}, {1, 1}, {0, 2}}},
		{id: 14, cells: []Coordinate{{0, 0}, {2, 0}, {0, 1}, {1, 1}, {2, 1}}},
		{id: 15, cells: []Coordinate{{0, 0}, {1, 0}, {2, 0}, {0, 1}, {0, 2}}},
		{id: 16, cells: []Coordinate{{0, 0}, {0, 1}, {1, 1}, {1, 2}, {2, 2}}},
		{id: 17, cells: []Coordinate{{0, 0}, {1, 0}, {1, 1}, {1, 2}, {2, 2}}},
		{id: 18, cells: []Coordinate{{1, 0}, {2, 0}, {0, 1}, {1, 1}, {1, 2}}},
		{id: 19, cells: []Coordinate{{1, 0}, {0, 1}, {1, 1}, {2, 1}, {1, 2}}},
		{id: 20, cells: []Coordinate{{0, 0}, {1, 0}, {2, 0}, {1, 1}, {1, 2}}},
	}
	for _, p := range pieces {
		p.normalize()
	}
	return pieces
}
