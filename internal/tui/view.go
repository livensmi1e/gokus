package tui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type viewData struct {
	turnLabel string

	leftPanel  string
	boardPanel string
	rightPanel string

	status string

	width  int
	height int
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
