package tui

import (
	"gokus/internal/room"

	tea "charm.land/bubbletea/v2"
)

type roomStateMsg struct {
	state room.State
}

type roomClosedMsg struct{}

func waitForRoomState(updates <-chan room.State) tea.Cmd {
	return func() tea.Msg {
		state, ok := <-updates
		if !ok {
			return roomClosedMsg{}
		}
		return roomStateMsg{state: state}
	}
}
