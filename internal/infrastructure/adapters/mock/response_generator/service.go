package mock

import "strings"

// Service provides a rotating set of canned responses.
// It implements the ResponsePort interface.
type Service struct {
	responses []string
	index     int
}

// NewService creates a Service with default responses.
func NewService() *Service {
	return &Service{
		responses: []string{
			"Entiendo. ¿Podrías darme más detalles?",
			"Interesante. Cuéntame más sobre eso.",
			"Estoy procesando tu mensaje...",
			"Gracias por compartir eso.",
			"¿Qué te gustaría hacer a continuación?",
		},
		index: 0,
	}
}

// Generate returns a rotating response based on the user input.
// Empty input returns a specific "no message" response.
func (s *Service) Generate(userInput string) string {
	if len(strings.TrimSpace(userInput)) == 0 {
		return "No has enviado ningún mensaje."
	}

	resp := s.responses[s.index]
	s.index = (s.index + 1) % len(s.responses)
	return resp
}
