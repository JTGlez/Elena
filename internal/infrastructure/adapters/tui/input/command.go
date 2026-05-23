package input

import (
	"strings"

	"elena/internal/core/domain"
	"elena/internal/core/usecases/identity"
)

// CommandHandler parses and executes slash-commands for the TUI.
type CommandHandler struct {
	identityUseCase *identity.IdentityUseCase
}

// NewCommandHandler creates a CommandHandler wired with IdentityUseCase.
func NewCommandHandler(identityUC *identity.IdentityUseCase) *CommandHandler {
	return &CommandHandler{
		identityUseCase: identityUC,
	}
}

// Parse checks if input starts with '/' and extracts the command token.
func (h *CommandHandler) Parse(input string) (cmd string, isCommand bool) {
	if len(input) > 0 && input[0] == '/' {
		parts := strings.SplitN(input, " ", 2)
		return parts[0], true
	}
	return "", false
}

// Execute dispatches a parsed command and returns the response string.
func (h *CommandHandler) Execute(cmd string, session *domain.Session) string {
	switch cmd {
	case "/exit":
		return "¡Hasta luego! Fue agradable hablar contigo."
	case "/reset":
		session.ClearMessages()
		return "Conversación reiniciada."
	case "/mood":
		mood := session.Mood().String()
		return "Estado actual: " + mood
	default:
		return "Comando no reconocido."
	}
}
