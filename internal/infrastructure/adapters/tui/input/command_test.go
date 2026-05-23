package input

import (
	"strings"
	"testing"

	"elena/internal/core/domain"
)

type parseScenario struct {
	name          string
	givenInput    string
	expectedCmd   string
	expectedIsCmd bool
}

type executeScenario struct {
	name             string
	givenCmd        string
	givenSession    *domain.Session
	expectedContains string
}

func TestCommandHandler_Parse_Command(t *testing.T) {
	scenario := getScenarioParseCommand()

	h := NewCommandHandler(nil)
	cmd, isCmd := h.Parse(scenario.givenInput)
	if cmd != scenario.expectedCmd {
		t.Fatalf("expected cmd %q, got %q", scenario.expectedCmd, cmd)
	}
	if isCmd != scenario.expectedIsCmd {
		t.Fatalf("expected isCommand %v, got %v", scenario.expectedIsCmd, isCmd)
	}
}

func TestCommandHandler_Parse_NonCommand(t *testing.T) {
	scenario := getScenarioParseNonCommand()

	h := NewCommandHandler(nil)
	cmd, isCmd := h.Parse(scenario.givenInput)
	if cmd != scenario.expectedCmd {
		t.Fatalf("expected cmd %q, got %q", scenario.expectedCmd, cmd)
	}
	if isCmd != scenario.expectedIsCmd {
		t.Fatalf("expected isCommand %v, got %v", scenario.expectedIsCmd, isCmd)
	}
}

func TestCommandHandler_Execute_Exit(t *testing.T) {
	scenario := getScenarioExecuteExit()

	h := NewCommandHandler(nil)
	result := h.Execute(scenario.givenCmd, scenario.givenSession)

	if !strings.Contains(result, scenario.expectedContains) {
		t.Fatalf("expected output to contain %q, got %q", scenario.expectedContains, result)
	}
}

func TestCommandHandler_Execute_Reset(t *testing.T) {
	scenario := getScenarioExecuteReset()

	h := NewCommandHandler(nil)
	result := h.Execute(scenario.givenCmd, scenario.givenSession)

	if !strings.Contains(result, scenario.expectedContains) {
		t.Fatalf("expected output to contain %q, got %q", scenario.expectedContains, result)
	}
	if scenario.givenSession.MessageCount() != 0 {
		t.Fatalf("expected 0 messages after reset, got %d", scenario.givenSession.MessageCount())
	}
}

func TestCommandHandler_Execute_MoodIdle(t *testing.T) {
	scenario := getScenarioExecuteMoodIdle()

	h := NewCommandHandler(nil)
	result := h.Execute(scenario.givenCmd, scenario.givenSession)

	if !strings.Contains(result, scenario.expectedContains) {
		t.Fatalf("expected output to contain %q, got %q", scenario.expectedContains, result)
	}
}

func TestCommandHandler_Execute_Unknown(t *testing.T) {
	scenario := getScenarioExecuteUnknown()

	h := NewCommandHandler(nil)
	result := h.Execute(scenario.givenCmd, scenario.givenSession)

	if !strings.Contains(result, scenario.expectedContains) {
		t.Fatalf("expected output to contain %q, got %q", scenario.expectedContains, result)
	}
}

func getScenarioParseCommand() parseScenario {
	return parseScenario{
		name:          "Given input starting with /, when Parse is called, then command token is returned with isCommand true",
		givenInput:    "/reset",
		expectedCmd:   "/reset",
		expectedIsCmd: true,
	}
}

func getScenarioParseNonCommand() parseScenario {
	return parseScenario{
		name:          "Given input without /, when Parse is called, then empty string is returned with isCommand false",
		givenInput:    "hola",
		expectedCmd:   "",
		expectedIsCmd: false,
	}
}

func getScenarioExecuteExit() executeScenario {
	return executeScenario{
		name:             "Given /exit command, when Execute is called, then goodbye message is returned",
		givenCmd:         "/exit",
		givenSession:     domain.NewSession("test-session"),
		expectedContains: "Hasta luego",
	}
}

func getScenarioExecuteReset() executeScenario {
	sess := domain.NewSession("test-session")
	sess.AddMessage(domain.NewMessage("", domain.AuthorUser, "mensaje"))
	return executeScenario{
		name:             "Given /reset command with messages, when Execute is called, then messages are cleared and confirmation returned",
		givenCmd:         "/reset",
		givenSession:     sess,
		expectedContains: "reiniciada",
	}
}

func getScenarioExecuteMoodIdle() executeScenario {
	sess := domain.NewSession("test-session")
	sess.SetMood(domain.MoodIdle)
	return executeScenario{
		name:             "Given /mood command with idle session, when Execute is called, then mood string is returned",
		givenCmd:         "/mood",
		givenSession:     sess,
		expectedContains: "idle",
	}
}

func getScenarioExecuteUnknown() executeScenario {
	return executeScenario{
		name:             "Given unknown command, when Execute is called, then not recognized message is returned",
		givenCmd:         "/unknown",
		givenSession:     domain.NewSession("test-session"),
		expectedContains: "no reconocido",
	}
}
