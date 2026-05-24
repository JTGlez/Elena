package output

// TUIDisplay implements output.DisplayPort.
// Notifications go through a callback into the BubbleTea model so they render
// in the TUI view, not on raw stdout.
type TUIDisplay struct {
	// Notify is called by Show(). Defaults to no-op. Set it after the model
	// is created: dp.Notify = m.SetNotification
	Notify func(string)
}

// NewTUIDisplay creates a TUIDisplay ready to have its Notify callback set.
func NewTUIDisplay() *TUIDisplay {
	return &TUIDisplay{}
}

// Show routes the notification content into the TUI via the Notify callback.
func (d *TUIDisplay) Show(content string) error {
	if d.Notify != nil {
		d.Notify(content)
	}
	return nil
}