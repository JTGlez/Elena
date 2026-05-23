package identity_test

import (
	"errors"
	"testing"

	"elena/internal/core/domain"
	"elena/internal/core/usecases/identity"

	"github.com/stretchr/testify/assert"
)

type mockDisplayPort struct {
	showCalled string
	showErr    error
}

func (m *mockDisplayPort) Show(content string) error {
	m.showCalled = content
	return m.showErr
}

type testScenario struct {
	name             string
	givenSessionMood string
	givenDisplayErr  error
	expectedMood     string
	expectedError    error
}

func TestIdentityUseCase_ShowMood_HappyPath(t *testing.T) {
	for _, scenario := range getSuccessScenarios() {
		t.Run(scenario.name, func(t *testing.T) {
			session := domain.NewSession("test-session")

			displayPort := &mockDisplayPort{
				showErr: scenario.givenDisplayErr,
			}

			uc := identity.New(identity.Input{
				DisplayPort: displayPort,
			})

			mood := uc.ShowMood(session)

			assert.Equal(t, scenario.expectedMood, mood)
			assert.Equal(t, scenario.givenSessionMood, displayPort.showCalled)
			assert.NoError(t, scenario.givenDisplayErr)
		})
	}
}

func TestIdentityUseCase_ShowMood_DisplayPortError(t *testing.T) {
	for _, scenario := range getDisplayErrorScenarios() {
		t.Run(scenario.name, func(t *testing.T) {
			session := domain.NewSession("test-session")

			displayPort := &mockDisplayPort{
				showErr: scenario.givenDisplayErr,
			}

			uc := identity.New(identity.Input{
				DisplayPort: displayPort,
			})

			mood := uc.ShowMood(session)

			assert.Equal(t, scenario.expectedMood, mood)
			assert.Equal(t, scenario.givenSessionMood, displayPort.showCalled)
		})
	}
}

func getSuccessScenarios() []testScenario {
	return []testScenario{
		{
			name:             "Given session with idle mood and working DisplayPort, when ShowMood is called, then the mood string is displayed and returned",
			givenSessionMood: "idle",
			givenDisplayErr:  nil,
			expectedMood:     "idle",
			expectedError:    nil,
		},
	}
}

func getDisplayErrorScenarios() []testScenario {
	displayErr := errors.New("display failed")

	return []testScenario{
		{
			name:             "Given DisplayPort returns an error, when ShowMood is called, then empty string is returned",
			givenSessionMood: "idle",
			givenDisplayErr:  displayErr,
			expectedMood:     "",
			expectedError:    displayErr,
		},
	}
}
