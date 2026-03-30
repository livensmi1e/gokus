package blokus

const DUO_BOARD_SIZE = 14

var STARTING_POINTS = []coordinate{
	{4, 4},
	{DUO_BOARD_SIZE - 5, DUO_BOARD_SIZE - 5},
}

type DuoGame struct {
	board         *board
	currentPlayer Occupant
	players       map[Occupant]*player
}

func NewDuoGame() *DuoGame {
	return &DuoGame{
		board:         newBoard(DUO_BOARD_SIZE),
		currentPlayer: Player1,
		players: map[Occupant]*player{
			Player1: {id: Player1, pieces: getPiecesAtStart()},
			Player2: {id: Player2, pieces: getPiecesAtStart()},
		},
	}
}

// Valid move
// Check for inbounds
// Empty space
// No edge to edge
// Must have corner to corner or player's first move
func (g *DuoGame) CanPlacePiece(id int, c coordinate) bool {
	player := g.players[g.currentPlayer]
	piece, ok := player.getPiece(id)
	if !ok {
		return false
	}
	touchedCorner := false
	coveredStartPoint := false
	for _, cell := range piece.cells {
		target := coordinate{c.x + cell.x, c.y + cell.y}
		// Check if piece can be placed at coordinate c
		if !g.board.isWithinBounds(target) || g.board.get(target) != Empty {
			return false
		}
		// If player first move
		if !player.started {
			// First move must be one of the starting points
			for _, sp := range STARTING_POINTS {
				if target == sp {
					coveredStartPoint = true
				}
			}
		}
		// Check for edge to edge contact and corner to corner
		if g.board.get(coordinate{target.x - 1, target.y}) == g.currentPlayer ||
			g.board.get(coordinate{target.x + 1, target.y}) == g.currentPlayer ||
			g.board.get(coordinate{target.x, target.y - 1}) == g.currentPlayer ||
			g.board.get(coordinate{target.x, target.y + 1}) == g.currentPlayer {
			return false
		}
		if g.board.get(coordinate{target.x - 1, target.y - 1}) == g.currentPlayer ||
			g.board.get(coordinate{target.x + 1, target.y - 1}) == g.currentPlayer ||
			g.board.get(coordinate{target.x - 1, target.y + 1}) == g.currentPlayer ||
			g.board.get(coordinate{target.x + 1, target.y + 1}) == g.currentPlayer {
			touchedCorner = true
		}
	}
	if !player.started {
		return coveredStartPoint
	}
	return touchedCorner
}

func (g *DuoGame) PlacePiece(id int, c coordinate) bool {
	if !g.CanPlacePiece(id, c) {
		return false
	}
	player := g.players[g.currentPlayer]
	piece, ok := player.takePiece(id)
	if !ok {
		return false
	}
	for _, cell := range piece.cells {
		target := coordinate{c.x + cell.x, c.y + cell.y}
		g.board.grid[target.y][target.x] = g.currentPlayer
	}
	if !player.started {
		player.started = true
	}
	g.turnPlayer()
	return true
}

func (g *DuoGame) IsOver() bool {
	for _, p := range g.players {
		if !p.stopped {
			return false
		}
	}
	return true
}

func (g *DuoGame) turnPlayer() {
	if g.currentPlayer == Player1 {
		g.currentPlayer = Player2
	} else {
		g.currentPlayer = Player1
	}
}
