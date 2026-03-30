package blokus

import "testing"

func TestBoard_Get_WithinBounds(t *testing.T) {
	b := newBoard(3)
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			got := b.get(Coordinate{x, y})
			if got != Empty {
				t.Errorf("expected Empty at (%d,%d), got %v", x, y, got)
			}
		}
	}
}

func TestBoard_Get_OutOfBounds(t *testing.T) {
	b := newBoard(3)
	tests := []Coordinate{
		{-1, 0},
		{0, -1},
		{3, 0},
		{0, 3},
		{100, 100},
		{-5, -2},
	}
	for _, c := range tests {
		got := b.get(c)
		if got != Empty {
			t.Errorf("expected Empty for out-of-bounds %v, got %v", c, got)
		}
	}
}

func TestBoard_IsWithinBounds(t *testing.T) {
	b := newBoard(3)
	tests := []struct {
		c    Coordinate
		want bool
	}{
		{Coordinate{0, 0}, true},
		{Coordinate{2, 2}, true},
		{Coordinate{1, 1}, true},
		{Coordinate{-1, 0}, false},
		{Coordinate{0, -1}, false},
		{Coordinate{3, 0}, false},
		{Coordinate{0, 3}, false},
		{Coordinate{100, 100}, false},
	}
	for _, tt := range tests {
		got := b.isWithinBounds(tt.c)
		if got != tt.want {
			t.Errorf("isWithinBounds(%v) = %v, want %v", tt.c, got, tt.want)
		}
	}
}

func TestBoard_SizeConsistency(t *testing.T) {
	b := newBoard(5)
	if len(b.grid) != 5 {
		t.Fatalf("expected grid height 5, got %d", len(b.grid))
	}
	for i := range b.grid {
		if len(b.grid[i]) != 5 {
			t.Errorf("expected row %d width 5, got %d", i, len(b.grid[i]))
		}
	}
}
