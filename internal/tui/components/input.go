// Package components provides reusable UI components for TUI screens.
package components

import (
	"flashcards/internal/tui/theme"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
)

// RenderLabeledInput renders a labeled text input with focus styling.
// It applies the appropriate style based on whether the input is focused.
func RenderLabeledInput(label string, input textinput.Model) string {
	labelRendered := theme.LabelStyle.Render(label)
	inputRendered := input.View()

	if input.Focused() {
		inputRendered = theme.InputFocusedStyle.Render(inputRendered)
	} else {
		inputRendered = theme.InputBlurredStyle.Render(inputRendered)
	}

	return fmt.Sprintf("%s\n%s", labelRendered, inputRendered)
}

// FormField represents a labeled input field in a form.
type FormField struct {
	Label string
	Input textinput.Model
}

// RenderFormFields renders multiple labeled inputs vertically.
// Each field is separated by a blank line for readability.
func RenderFormFields(fields ...FormField) string {
	var result string
	for i, field := range fields {
		if i > 0 {
			result += "\n\n"
		}
		result += RenderLabeledInput(field.Label, field.Input)
	}
	return result
}
