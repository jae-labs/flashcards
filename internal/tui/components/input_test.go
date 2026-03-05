package components

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
)

func TestRenderLabeledInput(t *testing.T) {
	tests := []struct {
		name     string
		label    string
		focused  bool
		contains []string
	}{
		{
			name:     "focused input",
			label:    "Username",
			focused:  true,
			contains: []string{"Username"},
		},
		{
			name:     "blurred input",
			label:    "Password",
			focused:  false,
			contains: []string{"Password"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := textinput.New()
			if tt.focused {
				input.Focus()
			} else {
				input.Blur()
			}

			result := RenderLabeledInput(tt.label, input)

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("RenderLabeledInput() = %q, want it to contain %q", result, expected)
				}
			}

			// Check that the result contains both label and input
			if !strings.Contains(result, "\n") {
				t.Errorf("RenderLabeledInput() should have label and input separated by newline")
			}
		})
	}
}

func TestRenderFormFields(t *testing.T) {
	tests := []struct {
		name     string
		fields   []FormField
		contains []string
	}{
		{
			name: "single field",
			fields: []FormField{
				{Label: "Field1", Input: textinput.New()},
			},
			contains: []string{"Field1"},
		},
		{
			name: "multiple fields",
			fields: []FormField{
				{Label: "Field1", Input: textinput.New()},
				{Label: "Field2", Input: textinput.New()},
			},
			contains: []string{"Field1", "Field2"},
		},
		{
			name:     "no fields",
			fields:   []FormField{},
			contains: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderFormFields(tt.fields...)

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("RenderFormFields() = %q, want it to contain %q", result, expected)
				}
			}

			// For multiple fields, verify they are separated
			if len(tt.fields) > 1 {
				if !strings.Contains(result, "\n\n") {
					t.Errorf("RenderFormFields() should separate fields with double newline")
				}
			}
		})
	}
}
