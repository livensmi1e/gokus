package tui

import (
	"context"
	"gokus/internal/blokus"
	"gokus/internal/room"
	"testing"
)

func TestWaitForRoomStateReturnsStateMessage(t *testing.T) {
	updates := make(chan room.State, 1)
	updates <- room.State{
		PlayerCount: 2,
	}
	msg := waitForRoomState(updates)()
	stateMsg, ok := msg.(roomStateMsg)
	if !ok {
		t.Fatalf("expected roomStateMsg, got %T", msg)
	}
	if stateMsg.state.PlayerCount != 2 {
		t.Fatalf(
			"expected player count 2, got %d",
			stateMsg.state.PlayerCount,
		)
	}
}

func TestWaitForRoomStateReturnsClosedMessage(t *testing.T) {
	updates := make(chan room.State)
	close(updates)
	msg := waitForRoomState(updates)()
	if _, ok := msg.(roomClosedMsg); !ok {
		t.Fatalf("expected roomClosedMsg, got %T", msg)
	}
}

func TestPlacePieceCmdReturnsResultMessage(t *testing.T) {
	r := room.New(context.Background())
	defer r.Close()
	client1, err := r.Join(context.Background())
	if err != nil {
		t.Fatalf("join player 1: %v", err)
	}
	if _, err := r.Join(context.Background()); err != nil {
		t.Fatalf("join player 2: %v", err)
	}
	cmd := placePieceCmd(
		client1,
		0,
		blokus.NewCoordinate(4, 4),
	)
	msg := cmd()
	result, ok := msg.(placeResultMsg)
	if !ok {
		t.Fatalf("expected placeResultMsg, got %T", msg)
	}
	if result.err != nil {
		t.Fatalf("place piece: %v", result.err)
	}
	if result.pieceId != 0 {
		t.Fatalf("expected piece ID 0, got %d", result.pieceId)
	}
}
