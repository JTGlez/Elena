package output

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"elena/internal/core/domain"
	"elena/internal/infrastructure/adapters/tui"
)

type Renderer struct{}

func NewRenderer() *Renderer {
	return &Renderer{}
}

// Render draws the split layout: avatar panel (left, fixed 25 cols) + chat panel (right).
// Takes an *Avatar so the renderer can display the current animated frame.
func (r *Renderer) Render(session *domain.Session, avatar *tui.Avatar, currentInput string) string {
	chat := r.renderChat(session, currentInput)
	avatarPanel := r.renderAvatarPanel(avatar)
	return lipgloss.JoinHorizontal(lipgloss.Top, avatarPanel, chat)
}

func (r *Renderer) renderAvatarPanel(avatar *tui.Avatar) string {
	frame := avatar.CurrentFrame()
	moodStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		Width(25).
		Align(lipgloss.Center)
	return moodStyle.Render(frame)
}

func (r *Renderer) renderChat(session *domain.Session, currentInput string) string {
	var b strings.Builder

	for _, msg := range session.Messages() {
		if msg.Author() == domain.AuthorUser {
			b.WriteString("> " + msg.Content() + "\n")
		} else {
			b.WriteString("Elena: " + msg.Content() + "\n")
		}
	}

	b.WriteString("\n> " + currentInput + "_")
	b.WriteString("\n\n/escribí tu mensaje y presioná Enter")
	b.WriteString("\n/exit  /reset  /mood (comandos)")

	chatStyle := lipgloss.NewStyle().
		Width(55).
		PaddingLeft(2)

	return chatStyle.Render(b.String())
}
