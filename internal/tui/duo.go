package tui

import (
	"gokus/internal/blokus"

	tea "charm.land/bubbletea/v2"
)

var _ tea.Model = &DouModel{}

type DouModel struct {
	game *blokus.DuoGame
}

func NewDuoModel() *DouModel {
	return &DouModel{
		game: blokus.NewDuoGame(),
	}
}

func (m *DouModel) Init() tea.Cmd {
	return nil
}

func (m *DouModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *DouModel) View() tea.View {
	s := "Ok go\n"
	return tea.NewView(s)
}
