package room

import (
	"context"
	"errors"
	"gokus/internal/blokus"
	"testing"
)

func TestRoomClose(t *testing.T) {
	ctx := context.Background()
	r := New(ctx)
	r.Close()
	select {
	case <-r.done:
		// done channel is closed
	default:
		t.Fatal("room should be closed already")
	}
}

func TestPlace(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	err := r.Place(
		context.Background(),
		0,
		blokus.NewCoordinate(4, 4),
	)
	if err != nil {
		t.Fatalf("expected valid move: %v", err)
	}
}

func TestPlaceInvalidMove(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	err := r.Place(
		context.Background(),
		0,
		blokus.NewCoordinate(0, 0),
	)
	if !errors.Is(err, ErrInvalidMove) {
		t.Fatalf("expected ErrInvalidMove, got %v", err)
	}
}

func TestPlaceAfterClose(t *testing.T) {
	r := New(context.Background())
	r.Close()
	err := r.Place(
		context.Background(),
		0,
		blokus.NewCoordinate(4, 4),
	)
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("expected ErrClosed, got %v", err)
	}
}

// This test is flaky because room maybe wating for request when request is sent
// Making 2 option is valid for select in Place
// Go will randomly choose one to continue
// Therefor not a good test because it does not test the expected contract of the API
// func TestPlaceWhenCallerContextAlreadyCancel(t *testing.T) {
// 	r := New(context.Background())
// 	defer r.Close()
// 	ctx, cancel := context.WithCancel(context.Background())
// 	cancel()
// 	err := r.Place(
// 		ctx,
// 		0,
// 		blokus.NewCoordinate(4, 4),
// 	)
// 	if !errors.Is(err, context.Canceled) {
// 		t.Fatalf("expected context canceled, got %v", err)
// 	}
// }
