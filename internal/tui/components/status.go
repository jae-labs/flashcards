// Package components provides reusable UI components for TUI screens.
package components

import (
	"flashcards/internal/tui/theme"
)

// StatusMessage holds error and success message state.
// Only one type of message (error or success) should be displayed at a time.
type StatusMessage struct {
	Error   string
	Success string
}

// Render returns a formatted status bar string.
// Error messages take precedence over success messages.
func (s *StatusMessage) Render() string {
	if s.Error != "" {
		return "\n" + theme.ErrorStyle.Render("✗ "+s.Error)
	}
	if s.Success != "" {
		return "\n" + theme.SuccessStyle.Render("✓ "+s.Success)
	}
	return ""
}

// SetError sets an error message and clears any success message.
func (s *StatusMessage) SetError(msg string) {
	s.Error = msg
	s.Success = ""
}

// SetSuccess sets a success message and clears any error message.
func (s *StatusMessage) SetSuccess(msg string) {
	s.Success = msg
	s.Error = ""
}

// Clear clears both error and success messages.
func (s *StatusMessage) Clear() {
	s.Error = ""
	s.Success = ""
}

// HasMessage returns true if either error or success message is set.
func (s *StatusMessage) HasMessage() bool {
	return s.Error != "" || s.Success != ""
}
