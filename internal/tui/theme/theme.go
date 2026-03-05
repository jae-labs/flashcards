// Package theme provides a centralized color palette and style definitions
// for all TUI components in the application. This ensures visual consistency
// and makes it easy to maintain and update the UI appearance.
package theme

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// Color palette - all colors used across the application
const (
	ColorPrimary     = "63"  // Purple/Cyan - borders, questions, titles
	ColorSuccess     = "34"  // Green - success messages
	ColorSuccessAlt  = "86"  // Light green - answers, selections, labels
	ColorError       = "196" // Red - errors
	ColorWarning     = "205" // Pink - table headers, warnings
	ColorInfo        = "244" // Gray - info text, unselected items
	ColorMuted       = "240" // Dark gray - blurred elements, borders
	ColorHighlight   = "229" // Yellow - focused text
	ColorHighlightBg = "57"  // Purple - focused backgrounds
	ColorCursor      = "212" // Pink - cursor indicator
)

// Text styles - common text formatting
var (
	// TitleStyle is used for screen titles and headings
	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorPrimary)).
			Bold(true)

	// QuestionStyle is used for displaying questions
	QuestionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(ColorPrimary))

	// AnswerStyle is used for displaying answers
	AnswerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSuccessAlt))

	// LabelStyle is used for form field labels
	LabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSuccessAlt)).
			Bold(true)
)

// Status styles - for feedback messages
var (
	// SuccessStyle is used for success messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSuccess)).
			Bold(true)

	// ErrorStyle is used for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorError)).
			Bold(true)

	// InfoStyle is used for informational messages
	InfoStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color(ColorInfo))

	// HelpStyle is used for help text
	HelpStyle = lipgloss.NewStyle().
			Faint(true)
)

// Interactive styles - for user interface elements
var (
	// SelectedStyle is used for selected items in lists
	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSuccessAlt))

	// UnselectedStyle is used for unselected items in lists
	UnselectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorInfo))

	// CursorStyle is used for the cursor indicator
	CursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorCursor))
)

// Input styles - for text input fields
var (
	// InputFocusedStyle is used for focused input fields
	InputFocusedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorHighlight)).
				Background(lipgloss.Color(ColorHighlightBg)).
				Padding(0, 1)

	// InputBlurredStyle is used for blurred/unfocused input fields
	InputBlurredStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorMuted)).
				Padding(0, 1)
)

// Checkbox styles - for checkbox elements
var (
	// CheckedStyle is used for checked checkboxes
	CheckedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorSuccess))

	// UncheckedStyle is used for unchecked checkboxes
	UncheckedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(ColorInfo))
)

// Table styles - for table components
var (
	// TableHeaderStyle is used for table headers
	TableHeaderStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color(ColorMuted)).
				BorderBottom(true).
				Bold(true).
				Foreground(lipgloss.Color(ColorWarning))

	// TableSelectedStyle is used for selected table rows
	TableSelectedStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(ColorHighlight)).
				Background(lipgloss.Color(ColorHighlightBg)).
				Bold(false)
)

// Layout constants - common dimensions and spacing
const (
	// MaxContentWidth is the maximum width for content frames
	MaxContentWidth = 80

	// MinTableHeight is the minimum height for table displays
	MinTableHeight = 5

	// MaxTableHeight is the maximum height for table displays
	MaxTableHeight = 25

	// DefaultPadding is the default padding for frames
	DefaultPadding = 1

	// FormPadding is the padding used for form layouts
	FormPadding = 2
)

// ApplyTableStyles applies the theme's table styles to a table model
func ApplyTableStyles(t table.Model) table.Model {
	s := table.DefaultStyles()
	s.Header = TableHeaderStyle
	s.Selected = TableSelectedStyle
	t.SetStyles(s)
	return t
}
