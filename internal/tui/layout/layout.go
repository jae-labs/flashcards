// Package layout provides helper functions for creating consistent layouts
// across all TUI screens. This includes frame creation, dimension calculations,
// and content positioning utilities.
package layout

import (
	"flashcards/internal/tui/theme"

	"github.com/charmbracelet/lipgloss"
)

// CalculateContentWidth returns a responsive width capped at the maximum content width.
// This ensures consistent frame widths across different terminal sizes.
func CalculateContentWidth(terminalWidth int) int {
	if terminalWidth > theme.MaxContentWidth {
		return theme.MaxContentWidth
	}
	return terminalWidth
}

// CreateFrame creates a standard bordered frame with the application's default styling.
// Additional options can be passed to customize the frame.
func CreateFrame(width int, opts ...FrameOption) lipgloss.Style {
	style := lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(theme.ColorPrimary)).
		Padding(theme.DefaultPadding, theme.DefaultPadding)

	for _, opt := range opts {
		style = opt(style)
	}

	return style
}

// FrameOption is a functional option for customizing frames.
type FrameOption func(lipgloss.Style) lipgloss.Style

// WithMaxHeight sets a maximum height for the frame.
func WithMaxHeight(h int) FrameOption {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.MaxHeight(h)
	}
}

// WithAlignment sets the horizontal and vertical alignment for the frame.
func WithAlignment(h, v lipgloss.Position) FrameOption {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.Align(h, v)
	}
}

// WithPadding sets custom padding for the frame.
func WithPadding(vertical, horizontal int) FrameOption {
	return func(s lipgloss.Style) lipgloss.Style {
		return s.Padding(vertical, horizontal)
	}
}

// CenterContent centers content on screen using lipgloss.Place.
func CenterContent(width, height int, content string) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// CalculateTableHeight computes a responsive table height based on terminal height.
// It accounts for various UI elements (borders, padding, help text, status messages).
func CalculateTableHeight(terminalHeight int) int {
	// Account for overhead: frame borders (2), padding (2), help (4), status (1), exit msg (3), centering (8)
	tableHeight := terminalHeight - 20
	if tableHeight < theme.MinTableHeight {
		tableHeight = theme.MinTableHeight
	}
	if tableHeight > theme.MaxTableHeight {
		tableHeight = theme.MaxTableHeight
	}
	return tableHeight
}

// CalculateMaxFrameHeight calculates the maximum frame height based on terminal height.
// It accounts for elements outside the frame (exit message, margins).
func CalculateMaxFrameHeight(terminalHeight int) int {
	// Account for: exit message outside frame (2), minimal margin for centering (4)
	maxHeight := terminalHeight - 6
	if maxHeight < 10 {
		maxHeight = 10
	}
	return maxHeight
}

// CalculateTableColumnWidths calculates responsive column widths for admin table.
// Returns widths for: ID, Question, Answer, and RevisitIn columns.
func CalculateTableColumnWidths(frameWidth int) (idWidth, questionWidth, answerWidth, revisitInWidth int) {
	// Account for: frame borders (2), frame padding left+right (4), table padding (2)
	availableWidth := frameWidth - 8

	// Default widths
	idWidth = 6
	revisitInWidth = 12

	// Adjust for very narrow screens
	if availableWidth < 50 {
		idWidth = 3
		revisitInWidth = 8
	}

	// Calculate widths for question and answer columns
	fixedWidth := idWidth + revisitInWidth + 4 // +4 for column spacing
	remainingWidth := availableWidth - fixedWidth
	if remainingWidth < 20 {
		remainingWidth = 20
	}
	questionWidth = int(float64(remainingWidth) * 0.55)
	answerWidth = remainingWidth - questionWidth

	return
}

// CreateBottomBar creates a centered status bar without borders for displaying
// left-aligned, center-aligned, and right-aligned text sections.
// Commonly used for displaying counters or status information at the bottom of views.
func CreateBottomBar(width int, left, center, right string) string {
	return lipgloss.NewStyle().
		Width(width - 4).
		Align(lipgloss.Center).
		Render(lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(10).Align(lipgloss.Left).Render(left),
			lipgloss.NewStyle().Width(width-24).Align(lipgloss.Center).Render(center),
			lipgloss.NewStyle().Width(10).Align(lipgloss.Right).Render(right),
		))
}
