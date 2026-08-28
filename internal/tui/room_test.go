package tui

import (
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
