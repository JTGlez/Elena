package tui_test

import (
	"testing"

	"elena/internal/core/domain"
	"elena/internal/infrastructure/adapters/tui"
)

type avatarScenario struct {
	name          string
	givenMood     domain.Mood
	givenTicks    int
	expectedFrame string
}

func TestAvatar_CurrentFrame_HappyPathIdle(t *testing.T) {
	scenario := getScenarioIdle()

	a := tui.NewAvatar(scenario.givenMood)
	got := a.CurrentFrame()

	if got != scenario.expectedFrame {
		t.Fatalf("expected frame %q, got %q", scenario.expectedFrame, got)
	}
}

func TestAvatar_CurrentFrame_TickAdvance(t *testing.T) {
	scenario := getScenarioTickAdvance()

	a := tui.NewAvatar(scenario.givenMood)
	for i := 0; i < scenario.givenTicks; i++ {
		a.Tick()
	}
	got := a.CurrentFrame()

	if got != scenario.expectedFrame {
		t.Fatalf("expected frame %q, got %q", scenario.expectedFrame, got)
	}
}

func TestAvatar_CurrentFrame_SetMoodReset(t *testing.T) {
	scenario := getScenarioSetMoodReset()

	a := tui.NewAvatar(scenario.givenMood)
	got := a.CurrentFrame()

	if got != scenario.expectedFrame {
		t.Fatalf("expected frame %q, got %q", scenario.expectedFrame, got)
	}
}

func TestAvatar_CurrentFrame_UnknownMood(t *testing.T) {
	scenario := getScenarioUnknownMood()

	a := tui.NewAvatar(scenario.givenMood)
	got := a.CurrentFrame()

	if got != scenario.expectedFrame {
		t.Fatalf("expected frame %q, got %q", scenario.expectedFrame, got)
	}
}

func getScenarioIdle() avatarScenario {
	return avatarScenario{
		name:          "Given idle mood, when CurrentFrame is called, then first idle frame is returned",
		givenMood:     domain.MoodIdle,
		givenTicks:    0,
		expectedFrame: "  /\\    /\\ \n ( ◕‿◕)\n  \\__/  \\__/",
	}
}

func getScenarioTickAdvance() avatarScenario {
	return avatarScenario{
		name:          "Given idle mood with 1 tick, when CurrentFrame is called, then second idle frame is returned",
		givenMood:     domain.MoodIdle,
		givenTicks:    1,
		expectedFrame: "  /\\    /\\ \n ( ◕ ‿ ◕)\n  \\__/  \\__/",
	}
}

func getScenarioSetMoodReset() avatarScenario {
	return avatarScenario{
		name:          "Given processing mood after SetMood, when CurrentFrame is called, then first processing frame is returned",
		givenMood:     domain.MoodProcessing,
		givenTicks:    0,
		expectedFrame: "  /\\    /\\ \n ( ◕ ○ ◕)\n  \\__/  \\__/",
	}
}

func getScenarioUnknownMood() avatarScenario {
	return avatarScenario{
		name:          "Given unknown mood value, when CurrentFrame is called, then empty string is returned",
		givenMood:     domain.Mood(99),
		givenTicks:    0,
		expectedFrame: "",
	}
}
