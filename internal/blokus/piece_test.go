package blokus

import (
	"reflect"
	"testing"
)

func TestNormalize_Table(t *testing.T) {
	tt := []struct {
		name  string
		input piece
		want  []Coordinate
	}{
		{
			name: "L shape shifted",
			input: piece{
				cells: []Coordinate{
					{2, 3}, {2, 4}, {3, 4},
				},
			},
			want: []Coordinate{
				{0, 0}, {0, 1}, {1, 1},
			},
		},
		{
			name: "line shifted",
			input: piece{
				cells: []Coordinate{
					{5, 5}, {6, 5}, {7, 5},
				},
			},
			want: []Coordinate{
				{0, 0}, {1, 0}, {2, 0},
			},
		},
		{
			name: "single cell",
			input: piece{
				cells: []Coordinate{
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
			p.normalize()

			if !reflect.DeepEqual(p.cells, tc.want) {
				t.Errorf("Normalize failed.\nGot: %v\nWant: %v", p.cells, tc.want)
			}
		})
	}
}

func TestRotate_Table(t *testing.T) {
	tt := []struct {
		name  string
		input piece
		want  []Coordinate
	}{
		{
			name: "L shape",
			input: piece{
				cells: []Coordinate{
					{0, 0}, {0, 1}, {1, 1},
				},
			},
			want: []Coordinate{
				{0, 0}, {0, 1}, {1, 0},
			},
		},
		{
			name: "line vertical -> horizontal",
			input: piece{
				cells: []Coordinate{
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
			p.rotate()
			p.normalize()
			if !reflect.DeepEqual(p.cells, tc.want) {
				t.Errorf("Rotate failed.\nGot: %v\nWant: %v", p.cells, tc.want)
			}
		})
	}
}

func TestFlip_Table(t *testing.T) {
	tt := []struct {
		name  string
		input piece
		want  []Coordinate
	}{
		{
			name: "basic L flip",
			input: piece{
				cells: []Coordinate{
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
			p.flip()
			p.normalize()
			if !reflect.DeepEqual(p.cells, tc.want) {
				t.Errorf("Flip failed.\nGot: %v\nWant: %v", p.cells, tc.want)
			}
		})
	}
}

func TestRotate4TimesReturnsOriginal_Table(t *testing.T) {
	tt := []struct {
		name     string
		input    piece
		expected piece
	}{
		{
			name: "rotate 4 times returns original L shape",
			input: piece{
				cells: []Coordinate{
					{0, 0},
					{1, 0},
					{1, 1},
				},
			},
			expected: piece{
				cells: []Coordinate{
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
				p.rotate()
				p.normalize()
			}
			if !reflect.DeepEqual(p, tc.expected) {
				t.Errorf("Rotate 4 times failed.\nGot: %v\nWant: %v", p, tc.expected)
			}
		})
	}
}
