package output

// ResponsePort generates Elena's responses to user input.
// The mock adapter provides a simple rotating response generator.
type ResponsePort interface {
	// Generate returns a response text for the given user input.
	Generate(userInput string) string
}
