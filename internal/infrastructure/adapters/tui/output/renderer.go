package output

import (
	"strings"

	"charm.land/lipgloss/v2"
	"elena/internal/core/domain"
	"elena/internal/infrastructure/adapters/tui"
)

// Renderer draws the split layout: avatar (25 cols) + chat panel.
// Call Render(session, avatar, input, notice) to get the full view string.
type Renderer struct{}

func NewRenderer() *Renderer {
	return &Renderer{}
}

// Render draws the split layout. The notice, if non-empty, is rendered INSIDE
// the chat panel (prepended at the top) to avoid lipgloss padding duplication
// that would otherwise duplicate it in the avatar panel.
func (r *Renderer) Render(session *domain.Session, avatar *tui.Avatar, currentInput, notice string) string {
	avatarPanel := r.renderAvatarPanel(avatar)
	chat := r.renderChat(session, currentInput, notice)
	return lipgloss.JoinHorizontal(lipgloss.Top, avatarPanel, chat)
}

func (r *Renderer) renderAvatarPanel(avatar *tui.Avatar) string {
	frame := avatar.CurrentFrame()
	moodStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("99")).
		Width(25).
		MaxWidth(25).
		Align(lipgloss.Center)
	return moodStyle.Render(frame)
}

func (r *Renderer) renderChat(session *domain.Session, currentInput, notice string) string {
	var b strings.Builder

	// Rule 5: prepend notice inside the panel so it does NOT cause
	// lipgloss padding duplication across the split layout.
	if notice != "" {
		b.WriteString(noticeStyle.Render(notice) + "\n")
	}

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
		MaxWidth(55).
		PaddingLeft(2)

	return chatStyle.Render(b.String())
}

var noticeStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("46")).
	Bold(true)