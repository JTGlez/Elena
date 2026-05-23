package output

import "fmt"

// TUIDisplay implements output.DisplayPort by printing to stdout.
// Used by the BubbleTea model to display content.
// In a full implementation this would render via the view system.
type TUIDisplay struct{}

// Show prints the content to stdout.
func (d *TUIDisplay) Show(content string) error {
	fmt.Println(content)
	return nil
}
