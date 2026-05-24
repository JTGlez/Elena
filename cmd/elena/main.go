package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"elena/internal/app"
	"elena/internal/infrastructure/adapters/tui/bubbletea"
)

func main() {
	a := app.Wire()

	// Build the model first so we can wire the displayPort callback.
	m := bubbletea.NewModel(
		a.Session(),
		a.ChatUseCase(),
		a.IdentityUseCase(),
		a.DisplayPort(),
	)

	// Bridge: displayPort notifications → model pendingNotice.
	a.DisplayPort().Notify = m.SetNotification

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("could not start: %v\n", err)
		os.Exit(1)
	}
}