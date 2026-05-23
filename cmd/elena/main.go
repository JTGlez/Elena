package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"elena/internal/app"
	"elena/internal/infrastructure/adapters/tui/bubbletea"
)

func main() {
	a := app.Wire()
	m := bubbletea.NewModel(a)
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("could not start: %v\n", err)
		os.Exit(1)
	}
}
