package chat

import (
	"context"
	"elena/internal/core/domain"
	"elena/internal/core/ports/output"
)

// Input holds the ports required by ChatUseCase.
type Input struct {
	ResponsePort output.ResponsePort
	DisplayPort  output.DisplayPort
}

// ChatUseCase orchestrates a single user-Elena message exchange.
type ChatUseCase struct {
	input Input
}

// New creates a new ChatUseCase with the given input ports.
func New(input Input) *ChatUseCase {
	return &ChatUseCase{input: input}
}

// Execute processes a user message: creates user and Elena messages,
// updates mood, generates response, and displays it.
func (c *ChatUseCase) Execute(ctx context.Context, session *domain.Session, userText string) (string, error) {
	userMsg := domain.NewMessage(generateID(), domain.AuthorUser, userText)
	session.AddMessage(userMsg)

	session.SetMood(domain.MoodProcessing)

	response := c.input.ResponsePort.Generate(userText)

	elenaMsg := domain.NewMessage(generateID(), domain.AuthorElena, response)
	session.AddMessage(elenaMsg)

	session.SetMood(domain.MoodIdle)

	if err := c.input.DisplayPort.Show(response); err != nil {
		return "", err
	}

	return response, nil
}

// generateID produces a message ID.
// TODO: replace with proper ID generator.
func generateID() string {
	return "placeholder-id"
}
