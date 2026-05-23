package domain

import (
	"testing"
	"time"
)

type sessionTestScenario struct {
	name                    string
	givenSession            *Session
	givenMessage            *Message
	givenNewMood            Mood
	givenID                 string
	expectedID              string
	expectedMood            Mood
	expectedCount           int
	expectedCountAfterModify int
	expectedMessages        []*Message
}

func TestSessionNew(t *testing.T) {
	scenario := getScenarioSessionNew()
	sess := NewSession(scenario.givenID)
	if sess.ID() != scenario.expectedID {
		t.Fatalf("expected ID %q, got %q", scenario.expectedID, sess.ID())
	}
	if sess.Mood() != scenario.expectedMood {
		t.Fatalf("expected Mood %v, got %v", scenario.expectedMood, sess.Mood())
	}
	if sess.MessageCount() != scenario.expectedCount {
		t.Fatalf("expected count %d, got %d", scenario.expectedCount, sess.MessageCount())
	}
	msgs := sess.Messages()
	if len(msgs) != scenario.expectedCount {
		t.Fatalf("expected Messages length %d, got %d", scenario.expectedCount, len(msgs))
	}
	if time.Since(sess.StartedAt()) > 5*time.Second {
		t.Fatalf("expected timestamp to be recent, but diff is %v", time.Since(sess.StartedAt()))
	}
}

func TestSessionGetID(t *testing.T) {
	scenario := getScenarioSessionGetID()
	got := scenario.givenSession.ID()
	if got != scenario.expectedID {
		t.Fatalf("expected %q, got %q", scenario.expectedID, got)
	}
}

func TestSessionGetMood(t *testing.T) {
	scenario := getScenarioSessionGetMood()
	got := scenario.givenSession.Mood()
	if got != scenario.expectedMood {
		t.Fatalf("expected %v, got %v", scenario.expectedMood, got)
	}
}

func TestSessionGetStartedAt(t *testing.T) {
	scenario := getScenarioSessionGetStartedAt()
	got := scenario.givenSession.StartedAt()
	if time.Since(got) > 5*time.Second {
		t.Fatalf("expected recent timestamp, got %v (diff %v)", got, time.Since(got))
	}
}

func TestSessionGetMessageCount(t *testing.T) {
	scenario := getScenarioSessionGetMessageCount()
	got := scenario.givenSession.MessageCount()
	if got != scenario.expectedCount {
		t.Fatalf("expected %d, got %d", scenario.expectedCount, got)
	}
}

func TestSessionAddMessage(t *testing.T) {
	scenario := getScenarioSessionAddMessage()
	scenario.givenSession.AddMessage(scenario.givenMessage)
	got := scenario.givenSession.MessageCount()
	if got != scenario.expectedCount {
		t.Fatalf("expected %d, got %d", scenario.expectedCount, got)
	}
	msgs := scenario.givenSession.Messages()
	if len(msgs) != scenario.expectedCount {
		t.Fatalf("expected Messages length %d, got %d", scenario.expectedCount, len(msgs))
	}
}

func TestSessionSetMood(t *testing.T) {
	scenario := getScenarioSessionSetMood()
	scenario.givenSession.SetMood(scenario.givenNewMood)
	got := scenario.givenSession.Mood()
	if got != scenario.expectedMood {
		t.Fatalf("expected %v, got %v", scenario.expectedMood, got)
	}
}

func TestSessionClearMessages(t *testing.T) {
	scenario := getScenarioSessionClearMessages()
	scenario.givenSession.ClearMessages()
	got := scenario.givenSession.MessageCount()
	if got != scenario.expectedCount {
		t.Fatalf("expected %d, got %d", scenario.expectedCount, got)
	}
}

func TestSessionMessagesDefensiveCopy(t *testing.T) {
	scenario := getScenarioSessionMessagesDefensiveCopy()
	scenario.givenSession.AddMessage(scenario.givenMessage)
	msgs := scenario.givenSession.Messages()
	msgs[0] = nil
	got := scenario.givenSession.MessageCount()
	if got != scenario.expectedCountAfterModify {
		t.Fatalf("expected %d, got %d", scenario.expectedCountAfterModify, got)
	}
}

func getScenarioSessionNew() sessionTestScenario {
	return sessionTestScenario{
		givenID:       "sess-001",
		expectedID:    "sess-001",
		expectedMood:  MoodIdle,
		expectedCount: 0,
	}
}

func getScenarioSessionGetID() sessionTestScenario {
	sess := NewSession("sess-abc")
	return sessionTestScenario{
		givenSession: sess,
		expectedID:   "sess-abc",
	}
}

func getScenarioSessionGetMood() sessionTestScenario {
	sess := NewSession("sess-001")
	return sessionTestScenario{
		givenSession: sess,
		expectedMood: MoodIdle,
	}
}

func getScenarioSessionGetStartedAt() sessionTestScenario {
	sess := NewSession("sess-001")
	return sessionTestScenario{
		givenSession: sess,
	}
}

func getScenarioSessionGetMessageCount() sessionTestScenario {
	sess := NewSession("sess-001")
	return sessionTestScenario{
		givenSession: sess,
		expectedCount: 0,
	}
}

func getScenarioSessionAddMessage() sessionTestScenario {
	sess := NewSession("sess-001")
	msg := NewMessage("msg-001", AuthorUser, "hello")
	return sessionTestScenario{
		givenSession:  sess,
		givenMessage:  msg,
		expectedCount: 1,
	}
}

func getScenarioSessionSetMood() sessionTestScenario {
	sess := NewSession("sess-001")
	return sessionTestScenario{
		givenSession: sess,
		givenNewMood: MoodProcessing,
		expectedMood: MoodProcessing,
	}
}

func getScenarioSessionClearMessages() sessionTestScenario {
	sess := NewSession("sess-001")
	msg := NewMessage("msg-001", AuthorUser, "to clear")
	sess.AddMessage(msg)
	return sessionTestScenario{
		givenSession:  sess,
		givenMessage:  msg,
		expectedCount: 0,
	}
}

func getScenarioSessionMessagesDefensiveCopy() sessionTestScenario {
	sess := NewSession("sess-001")
	msg := NewMessage("msg-001", AuthorUser, "hello")
	return sessionTestScenario{
		givenSession:           sess,
		givenMessage:           msg,
		expectedCountAfterModify: 1,
	}
}
