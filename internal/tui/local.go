package tui

import (
	"fmt"
	"gokus/internal/blokus"

	tea "charm.land/bubbletea/v2"
)

const (
	pieceCellWidth = 2 // width of a single cell of a piece in characters
	maxPieceCellsX = 5 // max width of any piece shape in cells
	pieceCardWidth = maxPieceCellsX * pieceCellWidth
	pieceCols      = 5
	pieceColGap    = 2

	HEADER = `
   _____       _ 
  / ____|     | |             
 | |  __  ___ | | ___   _ ___ 
 | | |_ |/ _ \| |/ / | | / __|
 | |__| | (_) |   <| |_| \__ \
  \_____|\___/|_|\_\\__,_|___/          
`

	FOOTER = "Keys: arrows move | tab next | shift+tab previous | enter place | r rotate | f flip | s skip | n reset | q quit"
)

var _ tea.Model = &LocalModel{}

type LocalModel struct {
	game *blokus.DuoGame

	styles Styles

	cursorX int
	cursorY int

	selectedPieceIdx int
	status           string

	width  int
	height int
}

func NewLocalModel() *LocalModel {
	return &LocalModel{
		game:   blokus.NewDuoGame(),
		styles: NewStyles(),
		status: "Welcome to Blokus Duo! Player 1 goes first.",
	}
}

func (m *LocalModel) Init() tea.Cmd {
	return nil
}

func (m *LocalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			m.tryPlaceSelectedPiece()
		case "r":
			m.rotateSelectedPiece()
		case "f":
			m.flipSelectedPiece()
		case "s":
			m.game.SkipTurn()
			m.selectedPieceIdx = 0
			m.status = "Turn skipped."
		case "n":
			m.game = blokus.NewDuoGame()
			m.cursorX, m.cursorY = 0, 0
			m.selectedPieceIdx = 0
			m.status = "New game started. Enjoy!"
		}
	}
	return m, nil
}

func (m *LocalModel) View() tea.View {
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

func (m *LocalModel) turnLabel() string {
	if m.game.IsOver() {
		return "Game over"
	}
	if m.game.CurrentPlayer() == blokus.Player1 {
		return "Player 1's Turn"
	}
	return "Player 2's Turn"
}

func (m *LocalModel) renderBoardPanel() string {
	return renderBoard(m.styles, boardViewData{
		board:      m.game.Board(),
		ghostCells: m.selectedPieceGhostCells(),
		cursor:     blokus.NewCoordinate(m.cursorX, m.cursorY),
	})
}

func (m *LocalModel) renderPlayerPanel(player blokus.Occupant) string {
	piecesLeft := m.game.PiecesLeft(player)
	isCurrentPlayer := m.game.CurrentPlayer() == player
	data := playerPanelViewData{
		player:          player,
		piecesLeft:      piecesLeft,
		squaresLeft:     m.game.Score(player),
		isCurrentPlayer: isCurrentPlayer,
	}
	if isCurrentPlayer && len(piecesLeft) > 0 {
		idx := min(m.selectedPieceIdx, len(piecesLeft)-1) // should we even need this min check here?
		data.activePieceID = piecesLeft[idx]
		data.hasActivePieceID = true
	}
	return renderPlayerPanel(m.styles, data)
}

func (m *LocalModel) cycleSelectedPiece(delta int) {
	pieces := m.game.PiecesLeft(m.game.CurrentPlayer())
	if len(pieces) == 0 {
		m.selectedPieceIdx = 0
		return
	}
	m.selectedPieceIdx = (m.selectedPieceIdx + delta + len(pieces)) % len(pieces)
}

func (m *LocalModel) tryPlaceSelectedPiece() {
	if m.game.IsOver() {
		m.status = "Game is over. Press r to restart."
		return
	}
	player := m.game.CurrentPlayer()
	pieces := m.game.PiecesLeft(player)
	if len(pieces) == 0 {
		m.status = "No pieces left. Press s to skip turn."
		return
	}
	m.selectedPieceIdx = min(m.selectedPieceIdx, len(pieces)-1)
	id := pieces[m.selectedPieceIdx]

	ok := m.game.PlacePiece(id, blokus.NewCoordinate(m.cursorX, m.cursorY))
	if ok {
		m.selectedPieceIdx = 0
		m.status = fmt.Sprintf("Player %d placed piece %d.", player, id)
		return
	}
	m.status = "Invalid move for this piece at current cursor."
}

func (m *LocalModel) rotateSelectedPiece() {
	pieces := m.game.PiecesLeft(m.game.CurrentPlayer())
	if len(pieces) == 0 {
		return
	}
	idx := min(m.selectedPieceIdx, len(pieces)-1)
	m.game.RotatePiece(pieces[idx])
}

func (m *LocalModel) flipSelectedPiece() {
	pieces := m.game.PiecesLeft(m.game.CurrentPlayer())
	if len(pieces) == 0 {
		return
	}
	idx := min(m.selectedPieceIdx, len(pieces)-1)
	m.game.FlipPiece(pieces[idx])
}

func (m *LocalModel) selectedPieceGhostCells() map[blokus.Coordinate]struct{} {
	ghostCells := make(map[blokus.Coordinate]struct{})
	if m.game.IsOver() {
		return ghostCells
	}

	pieces := m.game.PiecesLeft(m.game.CurrentPlayer())
	if len(pieces) == 0 {
		return ghostCells
	}

	idx := min(m.selectedPieceIdx, len(pieces)-1)
	shape := m.game.GetCurrentPieceShape(pieces[idx])
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
