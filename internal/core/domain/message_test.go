package domain

import (
	"testing"
	"time"
)

type messageTestScenario struct {
	name                string
	givenID             string
	givenAuthor         Author
	givenContent        string
	expectedMessage     *Message
	expectedID          string
	expectedAuthor      Author
	expectedContent     string
	expectedTimestampOK bool
}

func TestMessageNew(t *testing.T) {
	scenario := getScenarioMessageNew()
	msg := NewMessage(scenario.givenID, scenario.givenAuthor, scenario.givenContent)
	if msg == nil {
		t.Fatalf("expected non-nil message")
	}
	if msg.ID() != scenario.expectedID {
		t.Fatalf("expected ID %q, got %q", scenario.expectedID, msg.ID())
	}
	if msg.Author() != scenario.expectedAuthor {
		t.Fatalf("expected Author %q, got %q", scenario.expectedAuthor, msg.Author())
	}
	if msg.Content() != scenario.expectedContent {
		t.Fatalf("expected Content %q, got %q", scenario.expectedContent, msg.Content())
	}
	if !scenario.expectedTimestampOK {
		t.Fatalf("expected timestamp to be recent, but diff is %v", time.Since(msg.Timestamp()))
	}
}

func TestMessageGetID(t *testing.T) {
	scenario := getScenarioMessageGetID()
	got := scenario.expectedMessage.ID()
	if got != scenario.expectedID {
		t.Fatalf("expected %q, got %q", scenario.expectedID, got)
	}
}

func TestMessageGetAuthor(t *testing.T) {
	scenario := getScenarioMessageGetAuthor()
	got := scenario.expectedMessage.Author()
	if got != scenario.expectedAuthor {
		t.Fatalf("expected %q, got %q", scenario.expectedAuthor, got)
	}
}

func TestMessageGetContent(t *testing.T) {
	scenario := getScenarioMessageGetContent()
	got := scenario.expectedMessage.Content()
	if got != scenario.expectedContent {
		t.Fatalf("expected %q, got %q", scenario.expectedContent, got)
	}
}

func TestMessageGetTimestamp(t *testing.T) {
	scenario := getScenarioMessageGetTimestamp()
	got := scenario.expectedMessage.Timestamp()
	if time.Since(got).Abs() > 5*time.Second {
		t.Fatalf("expected recent timestamp, got %v (diff %v)", got, time.Since(got))
	}
}

func getScenarioMessageNew() messageTestScenario {
	return messageTestScenario{
		name:                "Given valid params, when NewMessage is called, then a message with all fields set is returned",
		givenID:             "msg-001",
		givenAuthor:         AuthorUser,
		givenContent:        "hello world",
		expectedID:          "msg-001",
		expectedAuthor:      AuthorUser,
		expectedContent:     "hello world",
		expectedTimestampOK: true,
	}
}

func getScenarioMessageGetID() messageTestScenario {
	msg := NewMessage("msg-042", AuthorElena, "test")
	return messageTestScenario{
		expectedMessage: msg,
		expectedID:      "msg-042",
	}
}

func getScenarioMessageGetAuthor() messageTestScenario {
	msg := NewMessage("msg-001", AuthorElena, "content")
	return messageTestScenario{
		expectedMessage: msg,
		expectedAuthor:  AuthorElena,
	}
}

func getScenarioMessageGetContent() messageTestScenario {
	msg := NewMessage("msg-001", AuthorUser, "bonjour")
	return messageTestScenario{
		expectedMessage: msg,
		expectedContent: "bonjour",
	}
}

func getScenarioMessageGetTimestamp() messageTestScenario {
	msg := NewMessage("msg-001", AuthorUser, "time test")
	now := time.Now()
	if now.Sub(msg.Timestamp()).Abs() > 5*time.Second {
		return messageTestScenario{
			expectedMessage:  msg,
			expectedTimestampOK: false,
		}
	}
	return messageTestScenario{
		expectedMessage:     msg,
		expectedTimestampOK: true,
	}
}
