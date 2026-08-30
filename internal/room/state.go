package room

import "gokus/internal/blokus"

// State is detached, point-in-time copy of a room's game state
// Mutate a State does not affect the room
type State struct {
	Board              [][]blokus.Occupant
	CurrentPlayer      blokus.Occupant
	PlayerCount        int
	PiecesLeft         map[blokus.Occupant][]int
	SquaresLeft        map[blokus.Occupant]int
	CurrentPieceShapes map[int][]blokus.Coordinate
	GameOver           bool
}
