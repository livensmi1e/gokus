package blokus

const DUO_BOARD_SIZE = 13

type DuoGame struct {
	board *board
}

func NewDuoGame() *DuoGame {
	return &DuoGame{
		board: newBoard(DUO_BOARD_SIZE),
	}
}
