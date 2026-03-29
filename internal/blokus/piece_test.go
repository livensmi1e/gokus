package blokus

import (
	"reflect"
	"testing"
)

func TestNormalize_Table(t *testing.T) {
	tt := []struct {
		name  string
		input Piece
		want  []Coordinate
	}{
		{
			name: "L shape shifted",
			input: Piece{
				Cells: []Coordinate{
					{2, 3}, {2, 4}, {3, 4},
				},
			},
			want: []Coordinate{
				{0, 0}, {0, 1}, {1, 1},
			},
		},
		{
			name: "line shifted",
			input: Piece{
				Cells: []Coordinate{
					{5, 5}, {6, 5}, {7, 5},
				},
			},
			want: []Coordinate{
				{0, 0}, {1, 0}, {2, 0},
			},
		},
		{
			name: "single cell",
			input: Piece{
				Cells: []Coordinate{
					{10, 10},
				},
			},
			want: []Coordinate{
				{0, 0},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.input
			p.Normalize()

			if !reflect.DeepEqual(p.Cells, tc.want) {
				t.Errorf("Normalize failed.\nGot: %v\nWant: %v", p.Cells, tc.want)
			}
		})
	}
}

func TestRotate_Table(t *testing.T) {
	tt := []struct {
		name  string
		input Piece
		want  []Coordinate
	}{
		{
			name: "L shape",
			input: Piece{
				Cells: []Coordinate{
					{0, 0}, {0, 1}, {1, 1},
				},
			},
			want: []Coordinate{
				{0, 0}, {0, 1}, {1, 0},
			},
		},
		{
			name: "line vertical -> horizontal",
			input: Piece{
				Cells: []Coordinate{
					{0, 0}, {0, 1}, {0, 2},
				},
			},
			want: []Coordinate{
				{0, 0}, {1, 0}, {2, 0},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.input
			p.Rotate()
			p.Normalize()
			if !reflect.DeepEqual(p.Cells, tc.want) {
				t.Errorf("Rotate failed.\nGot: %v\nWant: %v", p.Cells, tc.want)
			}
		})
	}
}

func TestFlip_Table(t *testing.T) {
	tt := []struct {
		name  string
		input Piece
		want  []Coordinate
	}{
		{
			name: "basic L flip",
			input: Piece{
				Cells: []Coordinate{
					{0, 0},
					{1, 0},
					{1, 1},
				},
			},
			want: []Coordinate{
				{0, 0},
				{0, 1},
				{1, 0},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.input
			p.Flip()
			p.Normalize()
			if !reflect.DeepEqual(p.Cells, tc.want) {
				t.Errorf("Flip failed.\nGot: %v\nWant: %v", p.Cells, tc.want)
			}
		})
	}
}

func TestRotate4TimesReturnsOriginal_Table(t *testing.T) {
	tt := []struct {
		name     string
		input    Piece
		expected Piece
	}{
		{
			name: "rotate 4 times returns original L shape",
			input: Piece{
				Cells: []Coordinate{
					{0, 0},
					{1, 0},
					{1, 1},
				},
			},
			expected: Piece{
				Cells: []Coordinate{
					{0, 0},
					{1, 0},
					{1, 1},
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			p := tc.input
			for i := 0; i < 4; i++ {
				p.Rotate()
				p.Normalize()
			}
			if !reflect.DeepEqual(p, tc.expected) {
				t.Errorf("Rotate 4 times failed.\nGot: %v\nWant: %v", p, tc.expected)
			}
		})
	}
}
