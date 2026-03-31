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
		status: "Welcome to Blokus Duo! Player 1 goes first.",
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

	body := lipgloss.JoinHorizontal(lipgloss.Center, left, "  ", center, "  ", right)
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
	ghostCells := m.selectedPieceGhostCells()
	for y := 0; y < len(board); y++ {
		var row strings.Builder
		for x := 0; x < len(board[y]); x++ {
			cell := board[y][x]
			style := m.cellStyle(cell)
			text := "  "
			if m.isStartingPoint(x, y) && cell == blokus.Empty {
				style = m.styles.Start
			}
			if _, ok := ghostCells[blokus.NewCoordinate(x, y)]; ok && cell == blokus.Empty {
				style = m.styles.Ghost
			}
			if _, ok := ghostCells[blokus.NewCoordinate(x, y)]; ok && cell != blokus.Empty {
				style = m.styles.Intersect
			}
			// if m.cursorX == x && m.cursorY == y {
			// 	style = m.styles.Cursor
			// 	text = "><"
			// }
			row.WriteString(style.Render(text))
		}
		lines = append(lines, row.String())
	}
	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return m.styles.Board.Render(content)
}

func (m *DouModel) renderPlayerPanel(player blokus.Occupant) string {
	header := "Player 1 (White)"
	if player == blokus.Player2 {
		header = "Player 2 (Black)"
	}

	pieceCount := m.game.Score(player)
	piecesLeft := m.game.PiecesLeft(player)
	scoreLine := fmt.Sprintf("Pieces: %d | Squares: %d", len(piecesLeft), pieceCount)

	activePieceID := -1
	if player == m.game.CurrentPlayer() && len(piecesLeft) > 0 {
		idx := min(m.selectedPieceIdx, len(piecesLeft)-1)
		activePieceID = piecesLeft[idx]
	}

	var cards []string
	for id := 0; id < 21; id++ {
		active := id == activePieceID
		used := !m.game.HasPiece(player, id)
		cards = append(cards, m.renderPieceCard(player, id, active, used))
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

func (m *DouModel) renderPieceCard(player blokus.Occupant, id int, active bool, used bool) string {
	shape := m.game.GetPieceShape(id)
	preview := m.renderPiece(player, shape, active, used)
	label := fmt.Sprintf("%02d", id)
	if active {
		label = m.styles.Cursor.Render(label)
	} else {
		label = m.styles.Subtle.Render(label)
	}
	card := lipgloss.JoinVertical(lipgloss.Center, preview, label)
	return lipgloss.NewStyle().Width(pieceCardWidth).Align(lipgloss.Center).Render(card)
}

func (m *DouModel) renderPiece(player blokus.Occupant, shape []blokus.Coordinate, active bool, used bool) string {
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
	if used {
		filled = m.styles.Used
	} else if player == blokus.Player2 {
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

func (m *DouModel) rotateSelectedPiece() {
	pieces := m.game.PiecesLeft(m.game.CurrentPlayer())
	if len(pieces) == 0 {
		return
	}
	idx := min(m.selectedPieceIdx, len(pieces)-1)
	m.game.RotatePiece(pieces[idx])
}

func (m *DouModel) flipSelectedPiece() {
	pieces := m.game.PiecesLeft(m.game.CurrentPlayer())
	if len(pieces) == 0 {
		return
	}
	idx := min(m.selectedPieceIdx, len(pieces)-1)
	m.game.FlipPiece(pieces[idx])
}

func (m *DouModel) selectedPieceGhostCells() map[blokus.Coordinate]struct{} {
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
