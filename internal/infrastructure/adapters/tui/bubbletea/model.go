package bubbletea

import (
	"elena/internal/app"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	app *app.App
}

func NewModel(a *app.App) tea.Model {
	return model{app: a}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	return "Elena MVP — Skeleton\nPress q to quit"
}
