package blokus

import "testing"

func newTestGame() *DuoGame { return NewDuoGame() }
func TestFirstMove_MustCoverStartingPoint(t *testing.T) {
	g := newTestGame()
	ok := g.PlacePiece(0, STARTING_POINTS[0])
	if !ok {
		t.Fatalf("expected first move to be valid on starting point")
	}
}
func TestFirstMove_InvalidStartPoint(t *testing.T) {
	g := newTestGame()
	ok := g.PlacePiece(0, Coordinate{0, 0})
	if ok {
		t.Fatalf("expected invalid first move to be rejected")
	}
}
func TestCannotOverlapPiece(t *testing.T) {
	g := newTestGame()
	start := STARTING_POINTS[0]
	if !g.PlacePiece(0, start) {
		t.Fatal("first placement failed")
	}
	if g.PlacePiece(0, start) {
		t.Fatal("expected overlap to fail")
	}
}
func TestTurnSwitching(t *testing.T) {
	g := newTestGame()
	start1 := STARTING_POINTS[0]
	start2 := STARTING_POINTS[1]
	if !g.PlacePiece(0, start1) {
		t.Fatal("player1 move failed")
	}
	if g.CurrentPlayer() != Player2 {
		t.Fatalf("expected Player2 turn")
	}
	if !g.PlacePiece(0, start2) {
		t.Fatal("player2 move failed")
	}
	if g.CurrentPlayer() != Player1 {
		t.Fatalf("expected Player1 turn")
	}
}
func TestEdgeTouchNotAllowed(t *testing.T) {
	g := newTestGame()
	start := STARTING_POINTS[0]
	if !g.PlacePiece(0, start) {
		t.Fatal("first move failed")
	}
	if g.PlacePiece(0, Coordinate{start.x + 1, start.y}) {
		t.Fatal("expected edge-touch to be invalid")
	}
}
func TestCannotReusePiece(t *testing.T) {
	g := newTestGame()
	if !g.PlacePiece(0, STARTING_POINTS[0]) {
		t.Fatal("p1 fail")
	}
	if !g.PlacePiece(0, STARTING_POINTS[1]) {
		t.Fatal("p2 fail")
	}
	if g.PlacePiece(0, Coordinate{5, 5}) {
		t.Fatal("reuse should fail")
	}
}
func TestCornerTouchRequiredAfterFirstMove(t *testing.T) {
	g := newTestGame()
	start := STARTING_POINTS[0]
	if !g.PlacePiece(0, start) {
		t.Fatal("first move failed")
	}
	if g.PlacePiece(1, Coordinate{start.x + 2, start.y}) {
		t.Fatal("should require corner touch")
	}
}
func TestCornerTouchValid(t *testing.T) {
	g := newTestGame()
	startOfPlayer1 := STARTING_POINTS[0]
	startOfPlayer2 := STARTING_POINTS[1]
	if !g.PlacePiece(0, startOfPlayer1) {
		t.Fatal("p1 first move failed")
	}
	if !g.PlacePiece(0, startOfPlayer2) {
		t.Fatal("p2 first move failed")
	}
	if !g.PlacePiece(1, Coordinate{startOfPlayer1.x + 1, startOfPlayer1.y + 1}) {
		t.Fatal("corner touch should be valid")
	}
}

func TestSkipTurn(t *testing.T) {
	g := newTestGame()
	p := g.CurrentPlayer()
	g.SkipTurn()
	if g.players[p].stopped != true {
		t.Fatal("player should be stopped")
	}
	if g.CurrentPlayer() == p {
		t.Fatal("turn should switch")
	}
}
func TestIsOver(t *testing.T) {
	g := newTestGame()
	g.SkipTurn()
	g.SkipTurn()
	if !g.IsOver() {
		t.Fatal("game should be over")
	}
}
func TestBoardIsImmutable(t *testing.T) {
	g := newTestGame()
	b := g.Board()
	b[0][0] = Player1
	if g.board.grid[0][0] == Player1 {
		t.Fatal("board should be immutable")
	}
}
func TestPiecesLeft(t *testing.T) {
	g := newTestGame()
	pieces := g.PiecesLeft(Player1)
	if len(pieces) != 21 {
		t.Fatal("should have 21 pieces at start")
	}
	g.PlacePiece(0, STARTING_POINTS[0])
	pieces = g.PiecesLeft(Player1)
	if len(pieces) != 20 {
		t.Fatal("should have 20 pieces after placement")
	}
}
func TestGetPieceShape(t *testing.T) {
	g := newTestGame()
	shape := g.GetPieceShape(0)
	if len(shape) != 1 {
		t.Fatal("monomino should have 1 cell")
	}
	shape[0].x = 999
	shape2 := g.GetPieceShape(0)
	if shape2[0].x == 999 {
		t.Fatal("shape should be immutable")
	}
}
func TestHasStarted(t *testing.T) {
	g := newTestGame()
	if g.HasStarted(Player1) {
		t.Fatal("should not start yet")
	}
	g.PlacePiece(0, STARTING_POINTS[0])
	if !g.HasStarted(Player1) {
		t.Fatal("should be started")
	}
}
func TestStartingPoints(t *testing.T) {
	g := newTestGame()
	points := g.StartingPoints()
	if len(points) != 2 {
		t.Fatal("should have 2 starting points")
	}
	points[0].x = 999
	points2 := g.StartingPoints()
	if points2[0].x == 999 {
		t.Fatal("should be immutable")
	}
}
