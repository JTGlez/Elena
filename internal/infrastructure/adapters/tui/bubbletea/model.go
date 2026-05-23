package bubbletea

import (
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	renderer        *output.Renderer
	quit            bool
}

// NewModel creates a new BubbleTea model wired with domain session and use cases.
func NewModel(session *domain.Session, chatUC *chat.ChatUseCase, identityUC *identity.IdentityUseCase) *model {
	return &model{
		session:         session,
		chatUseCase:     chatUC,
		identityUseCase: identityUC,
		avatar:          tui.NewAvatar(session.Mood()),
		cmdHandler:      input.NewCommandHandler(identityUC),
		currentInput:    "",
		renderer:        output.NewRenderer(),
		quit:            false,
	}
}

func (m *model) Init() tea.Cmd {
	// Ticker: 500ms per frame for avatar animation
	return tea.Every(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.avatar.Tick()
		return m, nil

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
			// Wire chatUseCase — executes in a goroutine via tea.Cmd
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
		} else if msg.Type == tea.KeyRunes {
			m.currentInput += msg.String()
		}

	case renderMsg:
		// Force re-render after async operations
		return m, nil
	}

	return m, nil
}

func (m *model) View() string {
	if m.quit {
		return lipgloss.NewStyle().Render("¡Hasta luego!\n")
	}
	return m.renderer.Render(m.session, m.avatar, m.currentInput)
}

func (m *model) handleCommand(cmd string) tea.Cmd {
	return func() tea.Msg {
		response := m.cmdHandler.Execute(cmd, m.session)
		if cmd == "/exit" {
			m.quit = true
		}
		m.currentInput = ""
		_ = response // response shown via DisplayPort in full impl; for MVP just execute
		return renderMsg{}
	}
}

type renderMsg struct{}

type tickMsg struct{}
