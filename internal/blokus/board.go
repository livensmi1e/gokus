package blokus

type Occupant int

const (
	Empty Occupant = iota
	Player1
	Player2
)

type board struct {
	size int
	grid [][]Occupant
}

func newBoard(size int) *board {
	grid := make([][]Occupant, size)
	for i := range grid {
		grid[i] = make([]Occupant, size)
	}
	return &board{size: size, grid: grid}
}

func (b *board) get(c coordinate) Occupant {
	return b.grid[c.y][c.x]
}
