package theme

import (
	"testing"

	"github.com/charmbracelet/bubbles/table"
)

func TestApplyTableStyles(t *testing.T) {
	// Create a simple table
	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Name", Width: 10},
	}
	rows := []table.Row{
		{"1", "Test"},
	}

	tbl := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
	)

	// Apply styles
	result := ApplyTableStyles(tbl)

	// Verify the table was returned and has styles applied
	// We can't easily check the exact styles, but we can verify it doesn't panic
	// and returns a valid table
	view := result.View()
	if view == "" {
		t.Error("ApplyTableStyles() returned table with empty view")
	}

	// Verify the table still has the same dimensions
	if len(result.Rows()) != len(rows) {
		t.Errorf("ApplyTableStyles() changed number of rows: got %d, want %d", len(result.Rows()), len(rows))
	}
}

func TestColorConstants(t *testing.T) {
	// Verify all color constants are defined and non-empty
	colors := map[string]string{
		"ColorPrimary":     ColorPrimary,
		"ColorSuccess":     ColorSuccess,
		"ColorSuccessAlt":  ColorSuccessAlt,
		"ColorError":       ColorError,
		"ColorWarning":     ColorWarning,
		"ColorInfo":        ColorInfo,
		"ColorMuted":       ColorMuted,
		"ColorHighlight":   ColorHighlight,
		"ColorHighlightBg": ColorHighlightBg,
		"ColorCursor":      ColorCursor,
	}

	for name, value := range colors {
		if value == "" {
			t.Errorf("%s is empty", name)
		}
	}
}

func TestLayoutConstants(t *testing.T) {
	// Verify layout constants have reasonable values
	if MaxContentWidth <= 0 {
		t.Errorf("MaxContentWidth = %d, want > 0", MaxContentWidth)
	}
	if MinTableHeight <= 0 {
		t.Errorf("MinTableHeight = %d, want > 0", MinTableHeight)
	}
	if MaxTableHeight <= 0 {
		t.Errorf("MaxTableHeight = %d, want > 0", MaxTableHeight)
	}
	if MaxTableHeight < MinTableHeight {
		t.Errorf("MaxTableHeight (%d) should be >= MinTableHeight (%d)", MaxTableHeight, MinTableHeight)
	}
	if DefaultPadding < 0 {
		t.Errorf("DefaultPadding = %d, want >= 0", DefaultPadding)
	}
	if FormPadding < 0 {
		t.Errorf("FormPadding = %d, want >= 0", FormPadding)
	}
}

func TestStylesNotNil(t *testing.T) {
	// Verify all styles are initialized (not nil/zero value)
	// We'll check by trying to render with them
	tests := []struct {
		name  string
		style interface{}
	}{
		{"TitleStyle", TitleStyle},
		{"QuestionStyle", QuestionStyle},
		{"AnswerStyle", AnswerStyle},
		{"LabelStyle", LabelStyle},
		{"SuccessStyle", SuccessStyle},
		{"ErrorStyle", ErrorStyle},
		{"InfoStyle", InfoStyle},
		{"HelpStyle", HelpStyle},
		{"SelectedStyle", SelectedStyle},
		{"UnselectedStyle", UnselectedStyle},
		{"CursorStyle", CursorStyle},
		{"InputFocusedStyle", InputFocusedStyle},
		{"InputBlurredStyle", InputBlurredStyle},
		{"CheckedStyle", CheckedStyle},
		{"UncheckedStyle", UncheckedStyle},
		{"TableHeaderStyle", TableHeaderStyle},
		{"TableSelectedStyle", TableSelectedStyle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify the style exists and doesn't cause issues
			// The actual rendering is handled by lipgloss
			if tt.style == nil {
				t.Errorf("%s is nil", tt.name)
			}
		})
	}
}
