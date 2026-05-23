package output

// DisplayPort is the output port for rendering messages in the UI.
// The UI adapter (BubbleTea) implements this to show content.
type DisplayPort interface {
	// Show renders the given content to the user.
	Show(content string) error
}
