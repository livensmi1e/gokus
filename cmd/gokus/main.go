package main

import (
	"gokus/internal/tui"
	"log"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	p := tea.NewProgram(tui.NewLocalModel())
	if _, err := p.Run(); err != nil {
		log.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
