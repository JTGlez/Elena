package domain

import "testing"

type moodTestScenario struct {
	name           string
	givenMood      Mood
	expectedString string
}

func TestMoodString_Idle(t *testing.T) {
	scenario := getScenarioMoodIdle()
	s := scenario.givenMood.String()
	if s != scenario.expectedString {
		t.Fatalf("expected %q, got %q", scenario.expectedString, s)
	}
}

func TestMoodString_Processing(t *testing.T) {
	scenario := getScenarioMoodProcessing()
	s := scenario.givenMood.String()
	if s != scenario.expectedString {
		t.Fatalf("expected %q, got %q", scenario.expectedString, s)
	}
}

func TestMoodString_Unknown(t *testing.T) {
	scenario := getScenarioMoodUnknown()
	s := scenario.givenMood.String()
	if s != scenario.expectedString {
		t.Fatalf("expected %q, got %q", scenario.expectedString, s)
	}
}

func getScenarioMoodIdle() moodTestScenario {
	return moodTestScenario{
		name:           "Given mood is idle, when String is called, then it returns idle",
		givenMood:      MoodIdle,
		expectedString: "idle",
	}
}

func getScenarioMoodProcessing() moodTestScenario {
	return moodTestScenario{
		name:           "Given mood is processing, when String is called, then it returns processing",
		givenMood:      MoodProcessing,
		expectedString: "processing",
	}
}

func getScenarioMoodUnknown() moodTestScenario {
	return moodTestScenario{
		name:           "Given mood is unknown, when String is called, then it returns unknown",
		givenMood:      Mood(99),
		expectedString: "unknown",
	}
}
