package bubbletea

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"elena/internal/core/domain"
	"elena/internal/core/usecases/chat"
	"elena/internal/core/usecases/identity"
	"elena/internal/infrastructure/adapters/tui"
	"elena/internal/infrastructure/adapters/tui/input"
	"elena/internal/infrastructure/adapters/tui/output"
)

// compile-time check that model satisfies tea.Model
var _ tea.Model = (*model)(nil)

type model struct {
	session         *domain.Session
	chatUseCase     *chat.ChatUseCase
	identityUseCase *identity.IdentityUseCase
	avatar          *tui.Avatar
	cmdHandler      *input.CommandHandler
	currentInput    string
	displayPort     *output.TUIDisplay
	renderer        *output.Renderer
	renderCount     int    // increments each View(); used to show notice exactly once
	pendingNotice   string
	quit            bool
}

// NewModel creates a new BubbleTea model wired with domain session, use cases,
// and display port (for notifications flowing into the TUI view).
func NewModel(
	session *domain.Session,
	chatUC *chat.ChatUseCase,
	identityUC *identity.IdentityUseCase,
	dp *output.TUIDisplay,
) *model {
	return &model{
		session:         session,
		chatUseCase:     chatUC,
		identityUseCase: identityUC,
		avatar:          tui.NewAvatar(session.Mood()),
		cmdHandler:      input.NewCommandHandler(identityUC),
		currentInput:    "",
		displayPort:     dp,
		renderer:        output.NewRenderer(),
		pendingNotice:   "",
		quit:            false,
	}
}

func (m *model) Init() tea.Cmd {
	// Rule 2: tea.Every auto-reschedules; tick handler returns nil.
	return tea.Every(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.avatar.Tick()
		return m, nil // tea.Every re-schedules itself

	case tea.KeyMsg:
		key := msg.String()
		if key == "q" || key == "ctrl+c" {
			return m, tea.Quit
		}
		if key == "enter" {
			if len(m.currentInput) == 0 {
				return m, nil
			}
			if cmd, ok := m.cmdHandler.Parse(m.currentInput); ok {
				return m, m.handleCommand(cmd)
			}
			// Rule 6: never call Update() from inside a Cmd.
			return m, func() tea.Msg {
				m.chatUseCase.Execute(nil, m.session, m.currentInput)
				m.avatar.SetMood(m.session.Mood())
				m.currentInput = ""
				return renderMsg{}
			}
		}
		if key == "backspace" {
			if len(m.currentInput) > 0 {
				m.currentInput = m.currentInput[:len(m.currentInput)-1]
			}
		} else if len(msg.Key().Text) > 0 {
			m.currentInput += msg.Key().Text
		}

	case renderMsg:
		return m, nil
	}
	return m, nil
}

func (m *model) View() tea.View {
	if m.quit {
		return tea.NewView(lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true).
			Render("¡Hasta luego!"))
	}
	// Show notice for exactly one render, then clear it.
	// The notice is captured locally, pendingNotice is cleared, and renderCount
	// is set to 1 so subsequent renders (e.g. ticker re-renders) suppress it.
	notice := ""
	if m.pendingNotice != "" {
		notice = m.pendingNotice
		m.pendingNotice = ""
		m.renderCount = 1
	}
	m.renderCount++
	v := tea.NewView(m.renderer.Render(m.session, m.avatar, m.currentInput, notice))
	v.AltScreen = true
	return v
}

// handleCommand executes a parsed slash command and returns the appropriate Cmd.
// The command response flows through displayPort → pendingNotice → View().
func (m *model) handleCommand(cmd string) tea.Cmd {
	return func() tea.Msg {
		response := m.cmdHandler.Execute(cmd, m.session)
		m.displayPort.Show(response) // sets m.pendingNotice via callback

		if cmd == "/exit" {
			m.quit = true
			// Rule 3: tea.QuitMsg inside func() tea.Msg, not tea.Quit.
			return func() tea.Msg {
				time.Sleep(2 * time.Second)
				return tea.QuitMsg{}
			}()
		}

		m.currentInput = ""
		return renderMsg{}
	}
}

// setNotification is the callback wired into displayPort so notifications
// land in the TUI view instead of stdout.
func (m *model) setNotification(notice string) {
	m.pendingNotice = notice
}

// Expose setNotification so TUIDisplay can call back into the model.
func (m *model) SetNotification(notice string) { m.setNotification(notice) }

type renderMsg struct{}
type tickMsg  struct{}