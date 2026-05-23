package output_test

import (
	"strings"
	"testing"

	"elena/internal/core/domain"
	"elena/internal/infrastructure/adapters/tui"
	"elena/internal/infrastructure/adapters/tui/output"
)

type renderScenario struct {
	name              string
	givenMessages     []messageInput
	givenInput        string
	expectedContains  []string
}

type messageInput struct {
	givenAuthor  domain.Author
	givenContent string
}

func TestRenderer_Render_SessionWithMessages(t *testing.T) {
	scenario := getScenarioWithMessages()

	sess := buildSession(scenario.givenMessages)
	givenAvatar := tui.NewAvatar(domain.MoodIdle)
	r := output.NewRenderer()
	result := r.Render(sess, givenAvatar, scenario.givenInput)

	for _, expected := range scenario.expectedContains {
		if !strings.Contains(result, expected) {
			t.Fatalf("expected output to contain %q, got:\n%s", expected, result)
		}
	}
}

func TestRenderer_Render_EmptySession(t *testing.T) {
	scenario := getScenarioEmptySession()

	sess := buildSession(scenario.givenMessages)
	givenAvatar := tui.NewAvatar(domain.MoodIdle)
	r := output.NewRenderer()
	result := r.Render(sess, givenAvatar, scenario.givenInput)

	for _, expected := range scenario.expectedContains {
		if !strings.Contains(result, expected) {
			t.Fatalf("expected output to contain %q, got:\n%s", expected, result)
		}
	}
}

func TestRenderer_Render_AvatarContainsEyes(t *testing.T) {
	scenario := getScenarioAvatarEyes()

	sess := buildSession(scenario.givenMessages)
	givenAvatar := tui.NewAvatar(domain.MoodIdle)
	r := output.NewRenderer()
	result := r.Render(sess, givenAvatar, scenario.givenInput)

	for _, expected := range scenario.expectedContains {
		if !strings.Contains(result, expected) {
			t.Fatalf("expected output to contain %q, got:\n%s", expected, result)
		}
	}
}

func buildSession(msgs []messageInput) *domain.Session {
	sess := domain.NewSession("test-session")
	for _, m := range msgs {
		sess.AddMessage(domain.NewMessage("", m.givenAuthor, m.givenContent))
	}
	return sess
}

func getScenarioWithMessages() renderScenario {
	return renderScenario{
		name: "Given session with user and elena messages, when Render is called, then both messages are rendered with correct prefixes",
		givenMessages: []messageInput{
			{
				givenAuthor:  domain.AuthorUser,
				givenContent: "Hola Elena",
			},
			{
				givenAuthor:  domain.AuthorElena,
				givenContent: "¡Hola! ¿Cómo estás?",
			},
		},
		givenInput:       "",
		expectedContains: []string{"> Hola Elena", "Elena: ¡Hola! ¿Cómo estás?"},
	}
}

func getScenarioEmptySession() renderScenario {
	return renderScenario{
		name: "Given empty session and input text, when Render is called, then prompt line with input and help text are rendered",
		givenMessages:  []messageInput{},
		givenInput:     "prueba",
		expectedContains: []string{"> prueba", "/escribí"},
	}
}

func getScenarioAvatarEyes() renderScenario {
	return renderScenario{
		name: "Given session with idle mood, when Render is called, then avatar contains eye character",
		givenMessages:  []messageInput{},
		givenInput:     "",
		expectedContains: []string{"◕"},
	}
}
