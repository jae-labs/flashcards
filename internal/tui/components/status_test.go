package components

import (
	"strings"
	"testing"
)

func TestStatusMessage_Render(t *testing.T) {
	tests := []struct {
		name     string
		error    string
		success  string
		contains string
	}{
		{
			name:     "error message",
			error:    "Something went wrong",
			success:  "",
			contains: "✗ Something went wrong",
		},
		{
			name:     "success message",
			error:    "",
			success:  "Operation successful",
			contains: "✓ Operation successful",
		},
		{
			name:     "error takes precedence",
			error:    "Error occurred",
			success:  "Success",
			contains: "✗ Error occurred",
		},
		{
			name:     "empty message",
			error:    "",
			success:  "",
			contains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StatusMessage{
				Error:   tt.error,
				Success: tt.success,
			}
			result := s.Render()

			if tt.contains == "" {
				if result != "" {
					t.Errorf("Render() = %q, want empty string", result)
				}
			} else {
				if !strings.Contains(result, tt.contains) {
					t.Errorf("Render() = %q, want it to contain %q", result, tt.contains)
				}
			}
		})
	}
}

func TestStatusMessage_SetError(t *testing.T) {
	s := &StatusMessage{
		Error:   "",
		Success: "Previous success",
	}

	s.SetError("New error")

	if s.Error != "New error" {
		t.Errorf("SetError() error = %q, want %q", s.Error, "New error")
	}
	if s.Success != "" {
		t.Errorf("SetError() success = %q, want empty string", s.Success)
	}
}

func TestStatusMessage_SetSuccess(t *testing.T) {
	s := &StatusMessage{
		Error:   "Previous error",
		Success: "",
	}

	s.SetSuccess("New success")

	if s.Success != "New success" {
		t.Errorf("SetSuccess() success = %q, want %q", s.Success, "New success")
	}
	if s.Error != "" {
		t.Errorf("SetSuccess() error = %q, want empty string", s.Error)
	}
}

func TestStatusMessage_Clear(t *testing.T) {
	s := &StatusMessage{
		Error:   "Some error",
		Success: "Some success",
	}

	s.Clear()

	if s.Error != "" {
		t.Errorf("Clear() error = %q, want empty string", s.Error)
	}
	if s.Success != "" {
		t.Errorf("Clear() success = %q, want empty string", s.Success)
	}
}

func TestStatusMessage_HasMessage(t *testing.T) {
	tests := []struct {
		name     string
		error    string
		success  string
		expected bool
	}{
		{
			name:     "has error",
			error:    "Error",
			success:  "",
			expected: true,
		},
		{
			name:     "has success",
			error:    "",
			success:  "Success",
			expected: true,
		},
		{
			name:     "has both",
			error:    "Error",
			success:  "Success",
			expected: true,
		},
		{
			name:     "has neither",
			error:    "",
			success:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StatusMessage{
				Error:   tt.error,
				Success: tt.success,
			}
			result := s.HasMessage()
			if result != tt.expected {
				t.Errorf("HasMessage() = %v, want %v", result, tt.expected)
			}
		})
	}
}
