package bubbletea

import (
	"strings"

	"elena/internal/app"
	"elena/internal/core/domain"
	"elena/internal/core/ports/output"

	tea "github.com/charmbracelet/bubbletea"
)

// Compile-time check that model satisfies tea.Model.
var _ tea.Model = (*model)(nil)

type model struct {
	app         *app.App
	currentInput string
	displayPort output.DisplayPort
}

// NewModel creates a new BubbleTea model with the given app and display port.
func NewModel(a *app.App, display output.DisplayPort) tea.Model {
	return model{
		app:         a,
		displayPort: display,
	}
}

// Init returns the initial command (none for this model).
func (m model) Init() tea.Cmd {
	return nil
}

// Update handles incoming messages.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		if key == "q" || key == "ctrl+c" {
			return m, tea.Quit
		}
		if key == "enter" {
			if len(m.currentInput) == 0 {
				return m, nil
			}
			// TODO: wire chatUseCase in PR-3 — for now just clear input
			m.currentInput = ""
			return m, nil
		}
		if key == "backspace" {
			if len(m.currentInput) > 0 {
				m.currentInput = m.currentInput[:len(m.currentInput)-1]
			}
		} else if msg.Type == tea.KeyRunes {
			m.currentInput += msg.String()
		}
	}
	return m, nil
}

// View renders the current state of the TUI.
func (m model) View() string {
	session := m.app.Session()
	lines := []string{"Elena MVP — Skeleton", ""}

	msgs := session.Messages()
	for _, msg := range msgs {
		prefix := "> "
		if msg.Author() == domain.AuthorElena {
			prefix = "Elena: "
		}
		lines = append(lines, prefix+msg.Content())
	}

	lines = append(lines, "")
	lines = append(lines, "> "+m.currentInput+"_")
	lines = append(lines, "")
	lines = append(lines, "(q to quit)")

	return strings.Join(lines, "\n") + "\n"
}
