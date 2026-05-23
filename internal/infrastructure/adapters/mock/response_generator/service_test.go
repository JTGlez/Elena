package mock_test

import (
	"testing"

	responsegen "elena/internal/infrastructure/adapters/mock/response_generator"

	"github.com/stretchr/testify/assert"
)

type testScenario struct {
	name            string
	givenUserInput  string
	givenResponses  []string
	givenIterations int
	expectedResponse string
}

func TestService_Generate_HappyPath(t *testing.T) {
	for _, scenario := range getSuccessScenarios() {
		t.Run(scenario.name, func(t *testing.T) {
			svc := responsegen.NewService()

			resp := svc.Generate(scenario.givenUserInput)

			assert.Equal(t, scenario.expectedResponse, resp)
		})
	}
}

func TestService_Generate_EmptyInput(t *testing.T) {
	for _, scenario := range getEmptyInputScenarios() {
		t.Run(scenario.name, func(t *testing.T) {
			svc := responsegen.NewService()

			resp := svc.Generate(scenario.givenUserInput)

			assert.Equal(t, scenario.expectedResponse, resp)
		})
	}
}

func TestService_Generate_RotationWraps(t *testing.T) {
	for _, scenario := range getRotationScenarios() {
		t.Run(scenario.name, func(t *testing.T) {
			svc := responsegen.NewService()

			var lastResp string
			for i := 0; i < scenario.givenIterations; i++ {
				lastResp = svc.Generate(scenario.givenUserInput)
			}

			assert.Equal(t, scenario.expectedResponse, lastResp)
		})
	}
}

func getSuccessScenarios() []testScenario {
	return []testScenario{
		{
			name:             "Given normal user input, when Generate is called, then first canned response is returned",
			givenUserInput:   "hello",
			givenIterations:  1,
			expectedResponse: "Entiendo. ¿Podrías darme más detalles?",
		},
	}
}

func getEmptyInputScenarios() []testScenario {
	return []testScenario{
		{
			name:             "Given empty user input, when Generate is called, then the no message response is returned",
			givenUserInput:   "",
			givenIterations:  1,
			expectedResponse: "No has enviado ningún mensaje.",
		},
	}
}

func getRotationScenarios() []testScenario {
	return []testScenario{
		{
			name:             "Given 5 successful Generate calls, when Generate is called again, then rotation wraps to first response",
			givenUserInput:   "hello",
			givenIterations:  6,
			expectedResponse: "Entiendo. ¿Podrías darme más detalles?",
		},
	}
}
