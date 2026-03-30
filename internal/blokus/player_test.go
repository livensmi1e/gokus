package blokus

import "testing"

func newTestPlayer() *player {
	p := &player{
		id:     1,
		pieces: make([]*piece, 21),
	}
	p.pieces[0] = &piece{}
	p.pieces[5] = &piece{}
	p.pieces[20] = &piece{}
	return p
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
		{1, false},
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
		{1, false, true},
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

func TestTakePiece(t *testing.T) {
	p := newTestPlayer()
	tt := []struct {
		first  int
		second int
	}{
		{5, 5},
	}
	for _, tc := range tt {
		pc, ok := p.takePiece(tc.first)
		if !ok || pc == nil {
			t.Fatalf("first take fail id=%d", tc.first)
		}
		if p.hasPiece(tc.first) {
			t.Fatalf("not removed id=%d", tc.first)
		}
		pc, ok = p.takePiece(tc.second)
		if ok || pc != nil {
			t.Fatalf("second take should fail id=%d", tc.second)
		}
	}
}
