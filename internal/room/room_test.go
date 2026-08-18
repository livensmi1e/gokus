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
//
// After adding pre-check context cancelled in Place, this test is not flaky anymore
func TestPlaceWhenCallerContextAlreadyCancel(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := r.Place(
		ctx,
		0,
		blokus.NewCoordinate(4, 4),
	)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context canceled, got %v", err)
	}
}

func TestJoinAssignsFirstPlayer(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	client, err := r.Join(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if client == nil {
		t.Fatal("expected a client, got nil")
	}
	if client.Player() != blokus.Player1 {
		t.Fatalf("expected Player1, got %v", client.Player())
	}
}

func TestJoinAssignsSecondPlayer(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	if _, err := r.Join(context.Background()); err != nil {
		t.Fatalf("join first player: %v", err)
	}
	client2, err := r.Join(context.Background())
	if err != nil {
		t.Fatalf("join second player: %v", err)
	}
	if client2 == nil {
		t.Fatal("expected a client, got nil")
	}
	if client2.Player() != blokus.Player2 {
		t.Fatalf("expected Player2, got %v", client2.Player())
	}
}

func TestJoinRejectsWhenRoomFull(t *testing.T) {
	r := New(context.Background())
	defer r.Close()
	if _, err := r.Join(context.Background()); err != nil {
		t.Fatalf("join first player: %v", err)
	}
	if _, err := r.Join(context.Background()); err != nil {
		t.Fatalf("join second player: %v", err)
	}
	client, err := r.Join(context.Background())
	if !errors.Is(err, ErrFull) {
		t.Fatalf("expected ErrFull, got: %v", err)
	}
	if client != nil {
		t.Fatal("expected nil client")
	}
}
