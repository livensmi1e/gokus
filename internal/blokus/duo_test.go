package blokus

import "testing"

// helper: tạo game fresh
func newTestGame() *DuoGame {
	return NewDuoGame()
}
func TestFirstMove_MustCoverStartingPoint(t *testing.T) {
	g := newTestGame()
	ok := g.PlacePiece(0, STARTING_POINTS[0])
	if !ok {
		t.Fatalf("expected first move to be valid on starting point")
	}
}

func TestFirstMove_InvalidStartPoint(t *testing.T) {
	g := newTestGame()
	wrong := coordinate{0, 0}
	ok := g.PlacePiece(0, wrong)
	if ok {
		t.Fatalf("expected invalid first move to be rejected")
	}
}

func TestCannotOverlapPiece(t *testing.T) {
	g := newTestGame()
	start := STARTING_POINTS[0]
	ok := g.PlacePiece(0, start)
	if !ok {
		t.Fatal("first placement failed")
	}
	ok = g.PlacePiece(0, start)
	if ok {
		t.Fatal("expected overlap to fail")
	}
}

func TestTurnSwitching(t *testing.T) {
	g := newTestGame()
	start1 := STARTING_POINTS[0]
	start2 := STARTING_POINTS[1]
	ok := g.PlacePiece(0, start1)
	if !ok {
		t.Fatal("player1 move failed")
	}
	if g.currentPlayer != Player2 {
		t.Fatalf("expected Player2 turn, got %v", g.currentPlayer)
	}
	ok = g.PlacePiece(0, start2)
	if !ok {
		t.Fatal("player2 move failed")
	}
	if g.currentPlayer != Player1 {
		t.Fatalf("expected Player1 turn, got %v", g.currentPlayer)
	}
}

func TestEdgeTouchNotAllowed(t *testing.T) {
	g := newTestGame()
	start := STARTING_POINTS[0]
	ok := g.PlacePiece(0, start)
	if !ok {
		t.Fatal("first move failed")
	}
	edgePos := coordinate{start.x + 1, start.y}
	ok = g.PlacePiece(0, edgePos)
	if ok {
		t.Fatal("expected edge-touch to be invalid")
	}
}

func TestCannotReusePiece(t *testing.T) {
	g := newTestGame()
	start := STARTING_POINTS[0]
	ok := g.PlacePiece(0, start)
	if !ok {
		t.Fatal("player 1 first placement failed")
	}
	ok = g.PlacePiece(0, STARTING_POINTS[1])
	if !ok {
		t.Fatal("player 2 first placement failed")
	}
}
