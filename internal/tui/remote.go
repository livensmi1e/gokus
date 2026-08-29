package tui

import (
	"context"
	"errors"
	"fmt"
	"gokus/internal/blokus"
	"gokus/internal/room"

	tea "charm.land/bubbletea/v2"
)

type roomStateMsg struct {
	state room.State
}

type roomClosedMsg struct{}

type placeResultMsg struct {
	pieceID int
	err     error
}

func waitForRoomState(updates <-chan room.State) tea.Cmd {
	return func() tea.Msg {
		state, ok := <-updates
		if !ok {
			return roomClosedMsg{}
		}
		return roomStateMsg{state: state}
	}
}

func placePieceCmd(client *room.Client, pieceID int, at blokus.Coordinate) tea.Cmd {
	return func() tea.Msg {
		err := client.Place(
			context.Background(),
			pieceID,
			at,
		)
		return placeResultMsg{
			pieceID: pieceID,
			err:     err,
		}
	}
}

var _ tea.Model = &RemoteModel{}

type RemoteModel struct {
	client *room.Client
	state  room.State

	styles Styles

	cursorX int
	cursorY int

	selectedPieceIdx int
	placing          bool
	status           string

	width  int
	height int
}

func NewRemoteModel(client *room.Client) *RemoteModel {
	return &RemoteModel{
		client: client,
		styles: NewStyles(),
		status: "Waiting for room state...",
	}
}

func (m *RemoteModel) Init() tea.Cmd {
	if m.client == nil {
		return nil
	}
	return waitForRoomState(m.client.Updates())
}

func (m *RemoteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "left":
			m.cursorX = max(0, m.cursorX-1)
		case "right":
			m.cursorX = min(blokus.DUO_BOARD_SIZE-1, m.cursorX+1)
		case "up":
			m.cursorY = max(0, m.cursorY-1)
		case "down":
			m.cursorY = min(blokus.DUO_BOARD_SIZE-1, m.cursorY+1)
		case "tab":
			m.cycleSelectedPiece(1)
		case "shift+tab":
			m.cycleSelectedPiece(-1)
		case "enter":
			if m.placing {
				return m, nil
			}
			pieceID, ok := m.selectedPieceID()
			if !ok {
				m.status = "No pieces left."
				return m, nil
			}
			at := blokus.NewCoordinate(m.cursorX, m.cursorY)
			m.placing = true
			m.status = "Placing piece..."
			return m, placePieceCmd(m.client, pieceID, at)
			// case "r":
			// 	m.rotateSelectedPiece()
			// case "f":
			// 	m.flipSelectedPiece()
			// case "s":
			// 	m.game.SkipTurn()
			// 	m.selectedPieceIDx = 0
			// 	m.status = "Turn skipped."
			// case "n":
			// 	m.game = blokus.NewDuoGame()
			// 	m.cursorX, m.cursorY = 0, 0
			// 	m.selectedPieceIDx = 0
			// 	m.status = "New game started. Enjoy!"
		}
	case roomStateMsg:
		previousPlayerCount := m.state.PlayerCount
		m.state = msg.state
		switch {
		case m.state.PlayerCount < 2:
			m.status = "Waiting for opponent..."
		case previousPlayerCount < 2:
			// only update status when room is ready in the first time
			m.status = "Game ready."
		}
		return m, waitForRoomState(m.client.Updates())
	case placeResultMsg:
		m.placing = false
		switch {
		case msg.err == nil:
			m.selectedPieceIdx = 0
			m.status = fmt.Sprintf(
				"Placed piece %d.",
				msg.pieceID,
			)
		case errors.Is(msg.err, room.ErrWaitingForOpponent):
			m.status = "Waiting for opponent."
		case errors.Is(msg.err, room.ErrOutOfTurn):
			m.status = "It is not your turn."
		case errors.Is(msg.err, room.ErrInvalidMove):
			m.status = "Invalid move at current position."
		case errors.Is(msg.err, room.ErrClosed):
			m.status = "Room closed."
		default:
			m.status = fmt.Sprintf(
				"Could not place piece: %v",
				msg.err,
			)
		}
		return m, nil
	case roomClosedMsg:
		m.status = "Room closed."
		return m, tea.Quit
	}
	return m, nil
}

func (m *RemoteModel) View() tea.View {
	return renderView(m.styles, viewData{
		turnLabel:  m.turnLabel(),
		leftPanel:  "",
		boardPanel: m.renderBoardPanel(),
		rightPanel: "",
		status:     m.status,
		width:      m.width,
		height:     m.height,
	})
}

func (m *RemoteModel) turnLabel() string {
	if m.state.PlayerCount < 2 {
		return "Waiting for opponent"
	}
	if m.state.CurrentPlayer == blokus.Player1 {
		return "Player 1's Turn"
	}
	return "Player 2's Turn"
}

func (m *RemoteModel) renderBoardPanel() string {
	return renderBoard(m.styles, boardViewData{
		board:      m.state.Board,
		ghostCells: nil,
		cursor:     blokus.NewCoordinate(m.cursorX, m.cursorY),
	})
}

func (m *RemoteModel) cycleSelectedPiece(delta int) {
	if m.client == nil {
		return
	}
	pieces := m.state.PiecesLeft[m.client.Player()]
	if len(pieces) == 0 {
		m.selectedPieceIdx = 0
		return
	}
	m.selectedPieceIdx = (m.selectedPieceIdx + delta + len(pieces)) % len(pieces)
}

func (m *RemoteModel) selectedPieceID() (int, bool) {
	if m.client == nil {
		return 0, false
	}
	pieces := m.state.PiecesLeft[m.client.Player()]
	if len(pieces) == 0 {
		return 0, false
	}
	m.selectedPieceIdx = min(m.selectedPieceIdx, len(pieces)-1)
	return pieces[m.selectedPieceIdx], true
}
