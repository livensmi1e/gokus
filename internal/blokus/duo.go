package blokus

const DUO_BOARD_SIZE = 14

var STARTING_POINTS = []Coordinate{
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

// CanPlacePiece checks if a piece can be placed at an particular point on board
// This is an immutable method
func (g *DuoGame) CanPlacePiece(id int, c Coordinate) bool {
	player := g.players[g.currentPlayer]
	if player.stopped {
		return false
	}
	if !player.hasPiece(id) {
		return false
	}
	piece, ok := player.getPiece(id)
	if !ok {
		return false
	}
	touchedCorner := false
	coveredStartPoint := false
	for _, cell := range piece.cells {
		target := Coordinate{c.x + cell.x, c.y + cell.y}
		// Check if piece can be placed at Coordinate c
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
		if g.board.get(Coordinate{target.x - 1, target.y}) == g.currentPlayer ||
			g.board.get(Coordinate{target.x + 1, target.y}) == g.currentPlayer ||
			g.board.get(Coordinate{target.x, target.y - 1}) == g.currentPlayer ||
			g.board.get(Coordinate{target.x, target.y + 1}) == g.currentPlayer {
			return false
		}
		if g.board.get(Coordinate{target.x - 1, target.y - 1}) == g.currentPlayer ||
			g.board.get(Coordinate{target.x + 1, target.y - 1}) == g.currentPlayer ||
			g.board.get(Coordinate{target.x - 1, target.y + 1}) == g.currentPlayer ||
			g.board.get(Coordinate{target.x + 1, target.y + 1}) == g.currentPlayer {
			touchedCorner = true
		}
	}
	if !player.started {
		return coveredStartPoint
	}
	return touchedCorner
}

// PlacePiece places a piece on a point in board
// This is mutable and changes the game state if the move is valid
func (g *DuoGame) PlacePiece(id int, c Coordinate) bool {
	if !g.CanPlacePiece(id, c) {
		return false
	}
	player := g.players[g.currentPlayer]
	piece, ok := player.getPiece(id)
	if !ok {
		return false
	}
	if !player.markPieceUsed(id) {
		return false
	}
	for _, cell := range piece.cells {
		target := Coordinate{c.x + cell.x, c.y + cell.y}
		g.board.grid[target.y][target.x] = g.currentPlayer
	}
	if !player.started {
		player.started = true
	}
	g.turnPlayer()
	return true
}

func (g *DuoGame) RotatePiece(id int) {
	player := g.players[g.currentPlayer]
	piece, ok := player.getPiece(id)
	if !ok {
		return
	}
	piece.rotate()
	piece.normalize()
}

func (g *DuoGame) FlipPiece(id int) {
	player := g.players[g.currentPlayer]
	piece, ok := player.getPiece(id)
	if !ok {
		return
	}
	piece.flip()
	piece.normalize()
}

func (g *DuoGame) SkipTurn() {
	g.players[g.currentPlayer].stopped = true
	g.turnPlayer()
}

func (g *DuoGame) IsOver() bool {
	for _, p := range g.players {
		if !p.stopped {
			return false
		}
	}
	return true
}

func (g *DuoGame) CurrentPlayer() Occupant {
	return g.currentPlayer
}

// Board return an immutable game board
func (g *DuoGame) Board() [][]Occupant {
	grid := make([][]Occupant, g.board.size)
	for i := range g.board.grid {
		grid[i] = append([]Occupant(nil), g.board.grid[i]...)
	}
	return grid
}

func (g *DuoGame) PiecesLeft(player Occupant) []int {
	p := g.players[player]
	var ids []int
	for i := range p.pieces {
		if p.hasPiece(i) {
			ids = append(ids, i)
		}
	}
	return ids
}

// PieceShape return the shape of a piece via coordinate of cells.
func PieceShape(id int) []Coordinate {
	pieces := getPiecesAtStart()
	pc := pieces[id]
	cells := make([]Coordinate, len(pc.cells))
	for i, c := range pc.cells {
		cells[i] = Coordinate{x: c.x, y: c.y}
	}
	return cells
}

// GetCurrentPieceShape return the current state of a piece (possibly rotated/flipped)
func (g *DuoGame) GetCurrentPieceShape(id int) []Coordinate {
	player := g.players[g.currentPlayer]
	pc, ok := player.getPiece(id)
	if !ok {
		return []Coordinate{}
	}
	cells := make([]Coordinate, len(pc.cells))
	for i, c := range pc.cells {
		cells[i] = Coordinate{x: c.x, y: c.y}
	}
	return cells
}

func (g *DuoGame) HasPiece(player Occupant, id int) bool {
	return g.players[player].hasPiece(id)
}

func (g *DuoGame) IsPieceUsed(player Occupant, id int) bool {
	return g.players[player].isPieceUsed(id)
}

func (g *DuoGame) GetPieceWidthAndHeight(player Occupant, id int) (int, int) {
	p := g.players[player]
	pc, ok := p.getPiece(id)
	if !ok {
		return 0, 0
	}
	return pc.getWidth(), pc.getHeight()
}

func (g *DuoGame) HasStarted(player Occupant) bool {
	return g.players[player].started
}

func (g *DuoGame) StartingPoints() []Coordinate {
	points := make([]Coordinate, len(STARTING_POINTS))
	for i, p := range STARTING_POINTS {
		points[i] = Coordinate{p.x, p.y}
	}
	return points
}

func (g *DuoGame) Score(player Occupant) int {
	p := g.players[player]
	score := 0
	for i := range p.pieces {
		if p.hasPiece(i) {
			score += len(p.pieces[i].cells)
		}
	}
	return score
}

func (g *DuoGame) turnPlayer() {
	nextPlayer := Player1
	if g.currentPlayer == Player1 {
		nextPlayer = Player2
	}
	if g.players[nextPlayer].stopped {
		return
	}
	g.currentPlayer = nextPlayer
}
