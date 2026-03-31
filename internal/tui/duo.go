package tui

import (
	"fmt"
	"gokus/internal/blokus"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	pieceCellWidth = 2
	maxPieceCellsX = 5
	pieceCardWidth = maxPieceCellsX * pieceCellWidth
	pieceCols      = 4
	pieceColGap    = 2
)

var _ tea.Model = &DouModel{}

type DouModel struct {
	game *blokus.DuoGame

	styles Styles

	cursorX int
	cursorY int

	selectedPieceIdx int
	status           string

	width  int
	height int
}

func NewDuoModel() *DouModel {
	return &DouModel{
		game:   blokus.NewDuoGame(),
		styles: NewStyles(),
		status: "Use arrows/hjkl to move, Tab to change piece, Enter to place.",
	}
}

func (m *DouModel) Init() tea.Cmd {
	return nil
}

func (m *DouModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "left", "h":
			m.cursorX = max(0, m.cursorX-1)
		case "right", "l":
			m.cursorX = min(blokus.DUO_BOARD_SIZE-1, m.cursorX+1)
		case "up", "k":
			m.cursorY = max(0, m.cursorY-1)
		case "down", "j":
			m.cursorY = min(blokus.DUO_BOARD_SIZE-1, m.cursorY+1)
		case "tab", "n":
			m.cycleSelectedPiece(1)
		case "shift+tab", "p":
			m.cycleSelectedPiece(-1)
		case "enter", " ":
			m.tryPlaceSelectedPiece()
		case "s":
			m.game.SkipTurn()
			m.selectedPieceIdx = 0
			m.status = "Turn skipped."
		case "r":
			m.game = blokus.NewDuoGame()
			m.cursorX, m.cursorY = 0, 0
			m.selectedPieceIdx = 0
			m.status = "New game started."
		}
	}
	return m, nil
}

func (m *DouModel) View() tea.View {
	headerText := normalizeHeaderArt(HEADER)
	header := lipgloss.JoinVertical(
		lipgloss.Center,
		m.styles.Title.Render(headerText),
		"",
		m.styles.Subtle.Render(m.turnLabel()),
	)

	left := m.renderPlayerPanel(blokus.Player1)
	center := m.renderBoardPanel()
	right := m.renderPlayerPanel(blokus.Player2)

	body := lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", center, "  ", right)
	bodyWidth := lipgloss.Width(body)
	header = lipgloss.NewStyle().Width(bodyWidth).Align(lipgloss.Center).Render(header)
	footer := m.styles.Subtle.Render(FOOTER)
	status := m.styles.Text.Render(m.status)

	ui := lipgloss.JoinVertical(lipgloss.Center, header, "", body, "", status, "", footer)
	rooted := m.styles.Root.Render(ui)
	if m.width > 0 && m.height > 0 {
		rooted = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, rooted)
	}
	return tea.NewView(rooted)
}

func (m *DouModel) turnLabel() string {
	if m.game.IsOver() {
		return "Game over"
	}
	if m.game.CurrentPlayer() == blokus.Player1 {
		return "Player 1's Turn"
	}
	return "Player 2's Turn"
}

func (m *DouModel) renderBoardPanel() string {
	var lines []string
	board := m.game.Board()
	for y := 0; y < len(board); y++ {
		var row strings.Builder
		for x := 0; x < len(board[y]); x++ {
			cell := board[y][x]
			style := m.cellStyle(cell)
			text := "  "
			if m.isStartingPoint(x, y) && cell == blokus.Empty {
				style = m.styles.Start
			}
			if m.cursorX == x && m.cursorY == y {
				style = m.styles.Cursor
				text = "><"
			}
			row.WriteString(style.Render(text))
		}
		lines = append(lines, row.String())
	}
	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return m.styles.Board.Render(content)
}

func (m *DouModel) renderPlayerPanel(player blokus.Occupant) string {
	pieceIDs := m.game.PiecesLeft(player)
	header := "Player 1 (White)"
	if player == blokus.Player2 {
		header = "Player 2 (Black)"
	}

	pieceCount := m.game.Score(player)
	scoreLine := fmt.Sprintf("Pieces: %d | Squares: %d", len(pieceIDs), pieceCount)

	var cards []string
	for idx, id := range pieceIDs {
		active := player == m.game.CurrentPlayer() && idx == m.selectedPieceIdx
		cards = append(cards, m.renderPieceCard(player, id, active))
	}
	if len(cards) == 0 {
		cards = []string{m.styles.Subtle.Render("No pieces left")}
	}

	piecesGrid := joinAsColumns(cards, pieceCols, pieceCardWidth, pieceColGap)

	titleStyle := m.styles.Title
	if player == blokus.Player2 {
		titleStyle = titleStyle.Foreground(lipgloss.Color("#CFCFD6"))
	}

	metaStyle := m.styles.Subtle
	if player == m.game.CurrentPlayer() {
		metaStyle = m.styles.Text
	}

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(header),
		metaStyle.Render(scoreLine),
		"",
		piecesGrid,
	)

	panelWidth := pieceCols*pieceCardWidth + (pieceCols-1)*pieceColGap + 4
	return m.styles.Panel.Width(panelWidth).Align(lipgloss.Center).Render(content)
}

func (m *DouModel) renderPieceCard(player blokus.Occupant, id int, active bool) string {
	shape := m.game.GetPieceShape(id)
	preview := m.renderPiecePreview(player, shape, active)
	label := fmt.Sprintf("%02d", id)
	if active {
		label = m.styles.Cursor.Render(label)
	} else {
		label = m.styles.Subtle.Render(label)
	}
	card := lipgloss.JoinVertical(lipgloss.Center, preview, label)
	return lipgloss.NewStyle().Width(pieceCardWidth).Align(lipgloss.Center).Render(card)
}

func (m *DouModel) renderPiecePreview(player blokus.Occupant, shape []blokus.Coordinate, active bool) string {
	if len(shape) == 0 {
		return ""
	}

	minX, minY := shape[0].X(), shape[0].Y()
	maxX, maxY := minX, minY
	for _, c := range shape {
		x, y := c.X(), c.Y()
		minX = min(minX, x)
		minY = min(minY, y)
		maxX = max(maxX, x)
		maxY = max(maxY, y)
	}

	w := maxX - minX + 1
	h := maxY - minY + 1
	filled := m.styles.Player1
	if player == blokus.Player2 {
		filled = m.styles.Player2
	}
	if active {
		filled = m.styles.Ghost
	}
	empty := lipgloss.NewStyle()

	localCoords := make(map[string]struct{}, len(shape))
	for _, c := range shape {
		localX := c.X() - minX
		localY := c.Y() - minY
		localCoords[fmt.Sprintf("%d:%d", localX, localY)] = struct{}{}
	}

	var lines []string
	for y := 0; y < h; y++ {
		var row strings.Builder
		for x := 0; x < w; x++ {
			key := fmt.Sprintf("%d:%d", x, y)
			if _, ok := localCoords[key]; ok {
				row.WriteString(filled.Render("  "))
				continue
			}
			row.WriteString(empty.Render("  "))
		}
		lines = append(lines, row.String())
	}

	block := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return lipgloss.NewStyle().Width(pieceCardWidth).Align(lipgloss.Center).Render(block)
}

func (m *DouModel) cellStyle(cell blokus.Occupant) lipgloss.Style {
	switch cell {
	case blokus.Player1:
		return m.styles.Player1
	case blokus.Player2:
		return m.styles.Player2
	default:
		return m.styles.Empty
	}
}

func (m *DouModel) isStartingPoint(x, y int) bool {
	for _, c := range m.game.StartingPoints() {
		if c.X() == x && c.Y() == y {
			return true
		}
	}
	return false
}

func (m *DouModel) cycleSelectedPiece(delta int) {
	pieces := m.game.PiecesLeft(m.game.CurrentPlayer())
	if len(pieces) == 0 {
		m.selectedPieceIdx = 0
		return
	}
	m.selectedPieceIdx = (m.selectedPieceIdx + delta + len(pieces)) % len(pieces)
}

func (m *DouModel) tryPlaceSelectedPiece() {
	if m.game.IsOver() {
		m.status = "Game is over. Press r to restart."
		return
	}

	pieces := m.game.PiecesLeft(m.game.CurrentPlayer())
	if len(pieces) == 0 {
		m.status = "No pieces left. Press s to skip turn."
		return
	}
	m.selectedPieceIdx = min(m.selectedPieceIdx, len(pieces)-1)
	id := pieces[m.selectedPieceIdx]

	ok := m.game.PlacePiece(id, blokus.NewCoordinate(m.cursorX, m.cursorY))
	if ok {
		m.selectedPieceIdx = 0
		m.status = fmt.Sprintf("Placed piece %d.", id)
		return
	}
	m.status = "Invalid move for this piece at current cursor."
}

func joinAsColumns(items []string, cols int, colWidth int, colGap int) string {
	if len(items) == 0 {
		return ""
	}
	rows := (len(items) + cols - 1) / cols
	gap := lipgloss.NewStyle().Width(colGap).Render("")
	var lines []string
	for r := 0; r < rows; r++ {
		rowItems := make([]string, 0, cols)
		for c := 0; c < cols; c++ {
			i := r*cols + c
			if i >= len(items) {
				rowItems = append(rowItems, lipgloss.NewStyle().Width(colWidth).Render(""))
				continue
			}
			rowItems = append(rowItems, lipgloss.NewStyle().Width(colWidth).Align(lipgloss.Center).Render(items[i]))
		}

		rowParts := make([]string, 0, len(rowItems)*2-1)
		for i, item := range rowItems {
			if i > 0 {
				rowParts = append(rowParts, gap)
			}
			rowParts = append(rowParts, item)
		}

		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Top, rowParts...))
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func normalizeHeaderArt(s string) string {
	trimmed := strings.Trim(s, "\n")
	lines := strings.Split(trimmed, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}

	minIndent := -1
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := len(line) - len(strings.TrimLeft(line, " \t"))
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	if minIndent <= 0 {
		return trimmed
	}

	for i, line := range lines {
		if len(line) >= minIndent {
			lines[i] = line[minIndent:]
		}
	}
	return strings.Join(lines, "\n")
}
