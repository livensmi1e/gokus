package tui

import (
	"fmt"
	"gokus/internal/blokus"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type viewData struct {
	turnLabel  string
	leftPanel  string
	boardPanel string
	rightPanel string
	status     string
	width      int
	height     int
}

func renderView(styles Styles, data viewData) tea.View {
	headerText := normalizeHeaderArt(HEADER)
	header := lipgloss.JoinVertical(
		lipgloss.Center,
		styles.Title.Render(headerText),
		"",
		styles.Subtle.Render(data.turnLabel),
	)
	body := lipgloss.JoinHorizontal(lipgloss.Center, data.leftPanel, "  ", data.boardPanel, "  ", data.rightPanel)
	bodyWidth := lipgloss.Width(body)
	header = lipgloss.NewStyle().Width(bodyWidth).Align(lipgloss.Center).Render(header)
	footer := styles.Subtle.Render(FOOTER)
	status := styles.Text.Render(data.status)
	ui := lipgloss.JoinVertical(lipgloss.Center, header, "", body, "", status, "", footer)
	rooted := styles.Root.Render(ui)
	if data.width > 0 && data.height > 0 {
		rooted = lipgloss.Place(data.width, data.height, lipgloss.Center, lipgloss.Center, rooted)
	}
	return tea.NewView(rooted)
}

type boardViewData struct {
	board      [][]blokus.Occupant
	ghostCells map[blokus.Coordinate]struct{}
	cursor     blokus.Coordinate
}

func renderBoard(styles Styles, data boardViewData) string {
	var lines []string
	board := data.board
	ghostCells := data.ghostCells
	for y := 0; y < len(board); y++ {
		var row strings.Builder
		for x := 0; x < len(board[y]); x++ {
			cell := board[y][x]
			style := cellStyle(styles, cell)
			text := "  "
			if isStartingPoint(blokus.STARTING_POINTS, x, y) && cell == blokus.Empty {
				style = styles.Start
			}
			if _, ok := ghostCells[blokus.NewCoordinate(x, y)]; ok && cell == blokus.Empty {
				style = styles.Ghost
			}
			if _, ok := ghostCells[blokus.NewCoordinate(x, y)]; ok && cell != blokus.Empty {
				style = styles.Intersect
			}
			// if data.cursor.X() == x && data.cursor.Y() == y {
			// 	style = styles.Cursor
			// 	text = "><"
			// }
			row.WriteString(style.Render(text))
		}
		lines = append(lines, row.String())
	}
	content := lipgloss.JoinVertical(lipgloss.Left, lines...)
	return styles.Board.Render(content)
}

type playerPanelViewData struct {
	player           blokus.Occupant
	isCurrentPlayer  bool
	piecesLeft       []int
	squaresLeft      int
	activePieceID    int
	hasActivePieceID bool
}

func renderPlayerPanel(styles Styles, data playerPanelViewData) string {
	header := "Player 1 (White)"
	if data.player == blokus.Player2 {
		header = "Player 2 (Black)"
	}
	scoreLine := fmt.Sprintf("Pieces: %d | Squares: %d", len(data.piecesLeft), data.squaresLeft)
	availablePieces := make(map[int]struct{}, len(data.piecesLeft))
	for _, pieceID := range data.piecesLeft {
		availablePieces[pieceID] = struct{}{}
	}
	var cards []string
	for pieceID := 0; pieceID < 21; pieceID++ {
		_, available := availablePieces[pieceID]
		active := data.hasActivePieceID && pieceID == data.activePieceID
		cards = append(cards, renderPieceCard(styles, data.player, pieceID, active, !available))
	}
	piecesGrid := joinAsColumns(cards, pieceCols, pieceCardWidth, pieceColGap)
	titleStyle := styles.Title
	if data.player == blokus.Player2 {
		titleStyle = titleStyle.Foreground(lipgloss.Color("#CFCFD6"))
	}
	metaStyle := styles.Subtle
	if data.isCurrentPlayer {
		metaStyle = styles.Text
	}
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render(header),
		metaStyle.Render(scoreLine),
		"",
		piecesGrid,
	)
	panelWidth := pieceCols*pieceCardWidth + (pieceCols-1)*pieceColGap + 4 // padding and border on both sides
	return styles.Panel.Width(panelWidth).Align(lipgloss.Center).Render(content)
}

func renderPieceCard(styles Styles, player blokus.Occupant, pieceID int, active, used bool) string {
	preview := renderPiece(styles, player, pieceID, active, used)
	label := fmt.Sprintf("%02d", pieceID+1)
	if active {
		label = styles.Cursor.Render(label)
	} else {
		label = styles.Subtle.Render(label)
	}
	card := lipgloss.JoinVertical(lipgloss.Center, preview, label)
	return lipgloss.NewStyle().Width(pieceCardWidth).Align(lipgloss.Center).Render(card)
}

func renderPiece(styles Styles, player blokus.Occupant, pieceID int, active bool, used bool) string {
	shape := blokus.PieceShape(pieceID)
	if len(shape) == 0 {
		return ""
	}
	w, h := pieceDimensions(shape)
	filled := styles.Player1
	if used {
		filled = styles.Used
	} else if player == blokus.Player2 {
		filled = styles.Player2
	}
	if active {
		filled = styles.Ghost
	}
	empty := lipgloss.NewStyle()
	coords := make(map[string]struct{}, len(shape))
	for _, c := range shape {
		coords[fmt.Sprintf("%d:%d", c.X(), c.Y())] = struct{}{}
	}
	var lines []string
	for y := 0; y < h; y++ {
		var row strings.Builder
		for x := 0; x < w; x++ {
			key := fmt.Sprintf("%d:%d", x, y)
			if _, ok := coords[key]; ok {
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

func pieceDimensions(
	shape []blokus.Coordinate,
) (width, height int) {
	for _, cell := range shape {
		width = max(width, cell.X()+1)
		height = max(height, cell.Y()+1)
	}
	return width, height
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

func cellStyle(styles Styles, cell blokus.Occupant) lipgloss.Style {
	switch cell {
	case blokus.Player1:
		return styles.Player1
	case blokus.Player2:
		return styles.Player2
	default:
		return styles.Empty
	}
}

func isStartingPoint(starts []blokus.Coordinate, x, y int) bool {
	for _, c := range starts {
		if c.X() == x && c.Y() == y {
			return true
		}
	}
	return false
}
