package tui

import (
	"context"
	"gokus/internal/blokus"
	"gokus/internal/room"

	tea "charm.land/bubbletea/v2"
)

type roomStateMsg struct {
	state room.State
}

type roomClosedMsg struct{}

type placeResultMsg struct {
	pieceId int
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

func placePieceCmd(client *room.Client, pieceId int, at blokus.Coordinate) tea.Cmd {
	return func() tea.Msg {
		err := client.Place(
			context.Background(),
			pieceId,
			at,
		)
		return placeResultMsg{
			pieceId: pieceId,
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
			m.selectedPieceId()
			pieceId, ok := m.selectedPieceId()
			if !ok {
				m.status = "No pieces left."
				return m, nil
			}
			at := blokus.NewCoordinate(m.cursorX, m.cursorY)
			m.status = "Placing piece..."
			return m, placePieceCmd(m.client, pieceId, at)
			// case "r":
			// 	m.rotateSelectedPiece()
			// case "f":
			// 	m.flipSelectedPiece()
			// case "s":
			// 	m.game.SkipTurn()
			// 	m.selectedPieceIdx = 0
			// 	m.status = "Turn skipped."
			// case "n":
			// 	m.game = blokus.NewDuoGame()
			// 	m.cursorX, m.cursorY = 0, 0
			// 	m.selectedPieceIdx = 0
			// 	m.status = "New game started. Enjoy!"
		}
	case roomStateMsg:
		m.state = msg.state
		if m.state.PlayerCount < 2 {
			m.status = "Waiting for opponent..."
		} else {
			m.status = "Game ready."
		}
		return m, waitForRoomState(m.client.Updates())
	case placeResultMsg:

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

func (m *RemoteModel) selectedPieceId() (int, bool) {
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
