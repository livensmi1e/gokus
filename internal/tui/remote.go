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

type rotateResultMsg struct {
	pieceID int
	err     error
}

type flipResultMsg struct {
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

func rotatePieceCmd(client *room.Client, pieceID int) tea.Cmd {
	return func() tea.Msg {
		err := client.Rotate(
			context.Background(),
			pieceID,
		)
		return rotateResultMsg{
			pieceID: pieceID,
			err:     err,
		}
	}
}

func flipPieceCmd(client *room.Client, pieceID int) tea.Cmd {
	return func() tea.Msg {
		err := client.Flip(
			context.Background(),
			pieceID,
		)
		return flipResultMsg{
			pieceID: pieceID,
			err:     err,
		}
	}
}

var _ tea.Model = &RemoteModel{}

type RemoteModel struct {
	client           *room.Client
	state            room.State
	styles           Styles
	cursorX          int
	cursorY          int
	selectedPieceIdx int
	cmdPending       bool
	status           string
	width            int
	height           int
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
			if m.cmdPending {
				return m, nil
			}
			pieceID, ok := m.selectedPieceID()
			if !ok {
				m.status = "No pieces left."
				return m, nil
			}
			at := blokus.NewCoordinate(m.cursorX, m.cursorY)
			m.cmdPending = true
			m.status = "Placing piece..."
			return m, placePieceCmd(m.client, pieceID, at)
		case "r":
			if m.cmdPending {
				return m, nil
			}
			pieceID, ok := m.selectedPieceID()
			if !ok {
				m.status = "No pieces left."
				return m, nil
			}
			m.cmdPending = true
			m.status = "Rotating piece..."
			return m, rotatePieceCmd(m.client, pieceID)
		case "f":
			if m.cmdPending {
				return m, nil
			}
			pieceID, ok := m.selectedPieceID()
			if !ok {
				m.status = "No pieces left."
				return m, nil
			}
			m.cmdPending = true
			m.status = "Flipping piece..."
			return m, flipPieceCmd(m.client, pieceID)
			// case "s":
			// 	m.game.SkipTurn()
			// 	m.selectedPieceIDx = 0
			// 	m.status = "Turn skipped."
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
		m.cmdPending = false
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
		case errors.Is(msg.err, room.ErrClosed):
			m.status = "Room closed."
		default:
			m.status = fmt.Sprintf(
				"Could not place piece: %v",
				msg.err,
			)
		}
		return m, nil
	case rotateResultMsg:
		m.cmdPending = false
		switch {
		case msg.err == nil:
			m.status = fmt.Sprintf(
				"Rotated piece %d.",
				msg.pieceID,
			)
		case errors.Is(msg.err, room.ErrWaitingForOpponent):
			m.status = "Waiting for opponent."
		case errors.Is(msg.err, room.ErrOutOfTurn):
			m.status = "It is not your turn."
		case errors.Is(msg.err, room.ErrPieceUnavailable):
			m.status = "Piece is no longer available."
		case errors.Is(msg.err, room.ErrClosed):
			m.status = "Room closed."
		default:
			m.status = fmt.Sprintf(
				"Could not rotate piece: %v",
				msg.err,
			)
		}
		return m, nil
	case flipResultMsg:
		m.cmdPending = false
		switch {
		case msg.err == nil:
			m.status = fmt.Sprintf(
				"Flipped piece %d.",
				msg.pieceID,
			)
		case errors.Is(msg.err, room.ErrWaitingForOpponent):
			m.status = "Waiting for opponent."
		case errors.Is(msg.err, room.ErrOutOfTurn):
			m.status = "It is not your turn."
		case errors.Is(msg.err, room.ErrPieceUnavailable):
			m.status = "Piece is no longer available."
		case errors.Is(msg.err, room.ErrInvalidMove):
			m.status = "Invalid move at current position."
		case errors.Is(msg.err, room.ErrClosed):
			m.status = "Room closed."
		default:
			m.status = fmt.Sprintf(
				"Could not flip piece: %v",
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
		leftPanel:  m.renderPlayerPanel(blokus.Player1),
		boardPanel: m.renderBoardPanel(),
		rightPanel: m.renderPlayerPanel(blokus.Player2),
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
		ghostCells: m.selectedPieceGhostCells(),
		cursor:     blokus.NewCoordinate(m.cursorX, m.cursorY),
	})
}

func (m *RemoteModel) renderPlayerPanel(player blokus.Occupant) string {
	piecesLeft := m.state.PiecesLeft[player]
	isCurrentPlayer := m.state.CurrentPlayer == player
	data := playerPanelViewData{
		player:          player,
		piecesLeft:      piecesLeft,
		squaresLeft:     m.state.SquaresLeft[player],
		isCurrentPlayer: isCurrentPlayer,
	}
	if m.client != nil && player == m.client.Player() && isCurrentPlayer && len(piecesLeft) > 0 {
		idx := min(m.selectedPieceIdx, len(piecesLeft)-1) // should we even need this min check here?
		data.activePieceID = piecesLeft[idx]
		data.hasActivePieceID = true
	}
	return renderPlayerPanel(m.styles, data)
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

func (m *RemoteModel) selectedPieceGhostCells() map[blokus.Coordinate]struct{} {
	ghostCells := make(map[blokus.Coordinate]struct{})
	if m.client == nil {
		return ghostCells
	}
	if m.client.Player() != m.state.CurrentPlayer {
		return ghostCells
	}
	pieceID, ok := m.selectedPieceID()
	if !ok {
		return ghostCells
	}
	shape, ok := m.state.CurrentPieceShapes[pieceID]
	if !ok {
		return ghostCells
	}
	for _, c := range shape {
		x := m.cursorX + c.X()
		y := m.cursorY + c.Y()
		if x < 0 || x >= blokus.DUO_BOARD_SIZE || y < 0 || y >= blokus.DUO_BOARD_SIZE {
			continue
		}
		ghostCells[blokus.NewCoordinate(x, y)] = struct{}{}
	}
	return ghostCells
}

func (m *RemoteModel) selectedPieceID() (int, bool) {
	if m.client == nil {
		return 0, false
	}
	pieces := m.state.PiecesLeft[m.client.Player()]
	if len(pieces) == 0 {
		return 0, false
	}
	idx := min(m.selectedPieceIdx, len(pieces)-1)
	return pieces[idx], true
}
