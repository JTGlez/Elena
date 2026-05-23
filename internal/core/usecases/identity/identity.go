package identity

import (
	"elena/internal/core/domain"
	"elena/internal/core/ports/output"
)

// Input holds the ports required by IdentityUseCase.
type Input struct {
	DisplayPort output.DisplayPort
}

// IdentityUseCase handles identity-related operations, such as mood display.
type IdentityUseCase struct {
	input Input
}

// New creates a new IdentityUseCase with the given input ports.
func New(input Input) *IdentityUseCase {
	return &IdentityUseCase{input: input}
}

// ShowMood renders the current session mood as a string via DisplayPort.
func (i *IdentityUseCase) ShowMood(session *domain.Session) string {
	mood := session.Mood().String()
	if err := i.input.DisplayPort.Show(mood); err != nil {
		return ""
	}
	return mood
}
