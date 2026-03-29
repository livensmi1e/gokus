package blokus

type Occupant int

const (
	Empty Occupant = iota
	Player1
	Player2
)

type Board struct {
	Size int
	Grid [][]Occupant
}

func NewBoard(size int) *Board {
	grid := make([][]Occupant, size)
	for i := range grid {
		grid[i] = make([]Occupant, size)
	}
	return &Board{Size: size, Grid: grid}
}

func (b Board) Get(c Coordinate) Occupant {
	return b.Grid[c.Y][c.X]
}
