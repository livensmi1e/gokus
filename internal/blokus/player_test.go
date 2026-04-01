package blokus

import "testing"

func newTestPlayer() *player {
	return &player{
		id:     1,
		pieces: getPiecesAtStart(),
	}
}

func TestHasPiece(t *testing.T) {
	p := newTestPlayer()
	tt := []struct {
		id  int
		exp bool
	}{
		{0, true},
		{5, true},
		{20, true},
		{1, true},
		{-1, false},
		{21, false},
	}
	for _, tc := range tt {
		if got := p.hasPiece(tc.id); got != tc.exp {
			t.Fatalf("id=%d exp=%v got=%v", tc.id, tc.exp, got)
		}
	}
}

func TestGetPiece(t *testing.T) {
	p := newTestPlayer()
	tt := []struct {
		id   int
		ok   bool
		nilp bool
	}{
		{0, true, false},
		{1, true, false},
		{-1, false, true},
	}
	for _, tc := range tt {
		pc, ok := p.getPiece(tc.id)
		if ok != tc.ok {
			t.Fatalf("id=%d expOk=%v got=%v", tc.id, tc.ok, ok)
		}
		if (pc == nil) != tc.nilp {
			t.Fatalf("id=%d nil=%v", tc.id, tc.nilp)
		}
	}
}

func TestMarkPieceUsed(t *testing.T) {
	p := newTestPlayer()
	tt := []struct {
		first  int
		second int
	}{
		{5, 5},
	}
	for _, tc := range tt {
		ok := p.markPieceUsed(tc.first)
		if !ok {
			t.Fatalf("first mark used fail id=%d", tc.first)
		}
		if p.hasPiece(tc.first) {
			t.Fatalf("should be unavailable after use id=%d", tc.first)
		}
		if !p.isPieceUsed(tc.first) {
			t.Fatalf("expected used flag true id=%d", tc.first)
		}
		ok = p.markPieceUsed(tc.second)
		if ok {
			t.Fatalf("second mark used should fail id=%d", tc.second)
		}
		if pc, ok := p.getPiece(tc.first); !ok || pc == nil {
			t.Fatalf("piece should still exist after use id=%d", tc.first)
		}
	}
}
