package chat_test

import (
	"context"
	"errors"
	"testing"

	"elena/internal/core/domain"
	"elena/internal/core/usecases/chat"

	"github.com/stretchr/testify/assert"
)

type mockResponsePort struct {
	generateResponse string
}

func (m *mockResponsePort) Generate(_ string) string {
	return m.generateResponse
}

type mockDisplayPort struct {
	showCalled string
	showErr    error
}

func (m *mockDisplayPort) Show(content string) error {
	m.showCalled = content
	return m.showErr
}

type testScenario struct {
	name              string
	givenUserText     string
	givenResponseText string
	givenDisplayErr   error
	expectedResponse  string
	expectedError     error
}

func TestChatUseCase_Execute_HappyPath(t *testing.T) {
	for _, scenario := range getSuccessScenarios() {
		t.Run(scenario.name, func(t *testing.T) {
			ctx := context.Background()

			responsePort := &mockResponsePort{
				generateResponse: scenario.givenResponseText,
			}

			displayPort := &mockDisplayPort{
				showErr: scenario.givenDisplayErr,
			}

			uc := chat.New(chat.Input{
				ResponsePort: responsePort,
				DisplayPort:  displayPort,
			})

			session := domain.NewSession("test-session")

			resp, err := uc.Execute(ctx, session, scenario.givenUserText)

			assert.NoError(t, err)
			assert.Equal(t, scenario.expectedResponse, resp)
			assert.Equal(t, scenario.givenResponseText, displayPort.showCalled)
			assert.Equal(t, 2, session.MessageCount())
			assert.Equal(t, domain.MoodIdle, session.Mood())
		})
	}
}

func TestChatUseCase_Execute_DisplayPortError(t *testing.T) {
	for _, scenario := range getDisplayErrorScenarios() {
		t.Run(scenario.name, func(t *testing.T) {
			ctx := context.Background()

			responsePort := &mockResponsePort{
				generateResponse: scenario.givenResponseText,
			}

			displayPort := &mockDisplayPort{
				showErr: scenario.givenDisplayErr,
			}

			uc := chat.New(chat.Input{
				ResponsePort: responsePort,
				DisplayPort:  displayPort,
			})

			session := domain.NewSession("test-session")

			resp, err := uc.Execute(ctx, session, scenario.givenUserText)

			assert.Equal(t, scenario.expectedError, err)
			assert.Equal(t, scenario.expectedResponse, resp)
		})
	}
}

func getSuccessScenarios() []testScenario {
	return []testScenario{
		{
			name:              "Given valid user text and working ports, when Execute is called, then response is returned with no error",
			givenUserText:     "hello",
			givenResponseText: "Elena's response",
			givenDisplayErr:   nil,
			expectedResponse:  "Elena's response",
			expectedError:     nil,
		},
	}
}

func getDisplayErrorScenarios() []testScenario {
	displayErr := errors.New("display failed")

	return []testScenario{
		{
			name:              "Given DisplayPort returns an error, when Execute is called, then the error is propagated",
			givenUserText:     "hello",
			givenResponseText: "Elena's response",
			givenDisplayErr:   displayErr,
			expectedResponse:  "",
			expectedError:     displayErr,
		},
	}
}
