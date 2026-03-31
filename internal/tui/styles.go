package tui

import "charm.land/lipgloss/v2"

type Styles struct {
	Root   lipgloss.Style
	Title  lipgloss.Style
	Subtle lipgloss.Style

	Board  lipgloss.Style
	Panel  lipgloss.Style
	Cursor lipgloss.Style

	Player1 lipgloss.Style
	Player2 lipgloss.Style

	Empty     lipgloss.Style
	Ghost     lipgloss.Style
	Start     lipgloss.Style
	Intersect lipgloss.Style
	Used      lipgloss.Style
	Text      lipgloss.Style
}

func NewStyles() Styles {
	return Styles{
		Root: lipgloss.NewStyle().
			Background(lipgloss.Color("#090B12")).
			Foreground(lipgloss.Color("#EAEAF0")).
			Padding(0, 0),
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F2F4FA")),
		Subtle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A6ADBD")),
		Board: lipgloss.NewStyle().
			Padding(2),
		Panel: lipgloss.NewStyle().
			Padding(0).
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("#232A3A")),
		Cursor: lipgloss.NewStyle().
			Background(lipgloss.Color("#F7D66A")).
			Foreground(lipgloss.Color("#1D1A14")).
			Bold(true),
		Player1: lipgloss.NewStyle().
			Background(lipgloss.Color("#EFEFEF")).
			Foreground(lipgloss.Color("#0F0F14")),
		Player2: lipgloss.NewStyle().
			Background(lipgloss.Color("#2D2A2D")).
			Foreground(lipgloss.Color("#F0EDF0")),
		Empty: lipgloss.NewStyle().
			Background(lipgloss.Color("#DCCFA4")).
			Foreground(lipgloss.Color("#A89767")),
		Start: lipgloss.NewStyle().
			Background(lipgloss.Color("#CBB67A")).
			Foreground(lipgloss.Color("#6E5F34")),
		Ghost: lipgloss.NewStyle().
			Background(lipgloss.Color("#bbc2ca")),
		Intersect: lipgloss.NewStyle().
			Background(lipgloss.Color("#e74141")),
		Used: lipgloss.NewStyle().
			Background(lipgloss.Color("#5B616B")).
			Foreground(lipgloss.Color("#D7DCE4")),
		Text: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EAEAF0")),
	}
}
