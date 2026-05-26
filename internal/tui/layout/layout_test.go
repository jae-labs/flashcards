package layout

import (
	"flashcards/internal/tui/theme"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestCalculateContentWidth(t *testing.T) {
	tests := []struct {
		name          string
		terminalWidth int
		expected      int
	}{
		{
			name:          "terminal wider than max",
			terminalWidth: 100,
			expected:      theme.MaxContentWidth,
		},
		{
			name:          "terminal narrower than max",
			terminalWidth: 60,
			expected:      60,
		},
		{
			name:          "terminal equals max",
			terminalWidth: theme.MaxContentWidth,
			expected:      theme.MaxContentWidth,
		},
		{
			name:          "very narrow terminal",
			terminalWidth: 20,
			expected:      20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateContentWidth(tt.terminalWidth)
			if result != tt.expected {
				t.Errorf("CalculateContentWidth(%d) = %d, want %d", tt.terminalWidth, result, tt.expected)
			}
		})
	}
}

func TestCreateFrame(t *testing.T) {
	tests := []struct {
		name  string
		width int
		opts  []FrameOption
	}{
		{
			name:  "default frame",
			width: 80,
			opts:  nil,
		},
		{
			name:  "frame with max height",
			width: 60,
			opts:  []FrameOption{WithMaxHeight(20)},
		},
		{
			name:  "frame with alignment",
			width: 70,
			opts:  []FrameOption{WithAlignment(lipgloss.Center, lipgloss.Center)},
		},
		{
			name:  "frame with custom padding",
			width: 50,
			opts:  []FrameOption{WithPadding(2, 3)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			style := CreateFrame(tt.width, tt.opts...)

			// Render some content to verify it works
			rendered := style.Render("Test content")
			if rendered == "" {
				t.Error("CreateFrame() style failed to render content")
			}

			// Verify it contains the test content
			if !strings.Contains(rendered, "Test content") {
				t.Error("CreateFrame() rendered output doesn't contain original content")
			}
		})
	}
}

func TestWithMaxHeight(t *testing.T) {
	opt := WithMaxHeight(10)
	style := lipgloss.NewStyle()
	result := opt(style)

	// Verify it renders without error
	rendered := result.Render("test")
	if rendered == "" {
		t.Error("WithMaxHeight() style failed to render")
	}
}

func TestWithAlignment(t *testing.T) {
	opt := WithAlignment(lipgloss.Center, lipgloss.Center)
	style := lipgloss.NewStyle()
	result := opt(style)

	// Verify it renders without error
	rendered := result.Render("test")
	if rendered == "" {
		t.Error("WithAlignment() style failed to render")
	}
}

func TestWithPadding(t *testing.T) {
	opt := WithPadding(2, 3)
	style := lipgloss.NewStyle()
	result := opt(style)

	// Verify it renders without error
	rendered := result.Render("test")
	if rendered == "" {
		t.Error("WithPadding() style failed to render")
	}
}

func TestCenterContent(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		height  int
		content string
	}{
		{
			name:    "simple content",
			width:   80,
			height:  24,
			content: "Hello World",
		},
		{
			name:    "empty content",
			width:   60,
			height:  20,
			content: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CenterContent(tt.width, tt.height, tt.content)

			// Result should not be empty
			if result == "" && tt.content != "" {
				t.Error("CenterContent() returned empty string for non-empty content")
			}

			// Should contain the original content
			if tt.content != "" && !strings.Contains(result, tt.content) {
				t.Errorf("CenterContent() = %q, want it to contain %q", result, tt.content)
			}
		})
	}
}

func TestCalculateTableHeight(t *testing.T) {
	tests := []struct {
		name           string
		terminalHeight int
		minExpected    int
		maxExpected    int
	}{
		{
			name:           "normal terminal",
			terminalHeight: 40,
			minExpected:    theme.MinTableHeight,
			maxExpected:    theme.MaxTableHeight,
		},
		{
			name:           "very tall terminal",
			terminalHeight: 100,
			minExpected:    theme.MinTableHeight,
			maxExpected:    theme.MaxTableHeight,
		},
		{
			name:           "short terminal",
			terminalHeight: 15,
			minExpected:    theme.MinTableHeight,
			maxExpected:    theme.MaxTableHeight,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateTableHeight(tt.terminalHeight)

			if result < tt.minExpected {
				t.Errorf("CalculateTableHeight(%d) = %d, want >= %d", tt.terminalHeight, result, tt.minExpected)
			}
			if result > tt.maxExpected {
				t.Errorf("CalculateTableHeight(%d) = %d, want <= %d", tt.terminalHeight, result, tt.maxExpected)
			}
		})
	}
}

func TestCalculateMaxFrameHeight(t *testing.T) {
	tests := []struct {
		name           string
		terminalHeight int
		minExpected    int
	}{
		{
			name:           "normal terminal",
			terminalHeight: 40,
			minExpected:    10,
		},
		{
			name:           "very tall terminal",
			terminalHeight: 100,
			minExpected:    10,
		},
		{
			name:           "short terminal",
			terminalHeight: 15,
			minExpected:    9,
		},
		{
			name:           "very short terminal",
			terminalHeight: 5,
			minExpected:    10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateMaxFrameHeight(tt.terminalHeight)

			if result < tt.minExpected {
				t.Errorf("CalculateMaxFrameHeight(%d) = %d, want >= %d", tt.terminalHeight, result, tt.minExpected)
			}
		})
	}
}

func TestCalculateTableColumnWidths(t *testing.T) {
	tests := []struct {
		name       string
		frameWidth int
	}{
		{
			name:       "normal width",
			frameWidth: 80,
		},
		{
			name:       "narrow width",
			frameWidth: 50,
		},
		{
			name:       "very narrow width",
			frameWidth: 30,
		},
		{
			name:       "wide width",
			frameWidth: 120,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idWidth, questionWidth, answerWidth, revisitInWidth := CalculateTableColumnWidths(tt.frameWidth)

			// All widths should be positive
			if idWidth <= 0 || questionWidth <= 0 || answerWidth <= 0 || revisitInWidth <= 0 {
				t.Errorf("CalculateTableColumnWidths(%d) returned non-positive width(s): id=%d, q=%d, a=%d, r=%d",
					tt.frameWidth, idWidth, questionWidth, answerWidth, revisitInWidth)
			}

			// Question width should be larger than answer width (55/45 split)
			if questionWidth < answerWidth {
				t.Errorf("CalculateTableColumnWidths(%d): questionWidth=%d should be >= answerWidth=%d",
					tt.frameWidth, questionWidth, answerWidth)
			}
		})
	}
}

func TestCreateBottomBar(t *testing.T) {
	tests := []struct {
		name   string
		width  int
		left   string
		center string
		right  string
	}{
		{
			name:   "with all sections",
			width:  80,
			left:   "Left",
			center: "Center",
			right:  "Right",
		},
		{
			name:   "with empty sections",
			width:  60,
			left:   "",
			center: "",
			right:  "",
		},
		{
			name:   "with partial content",
			width:  70,
			left:   "1/10",
			center: "",
			right:  "100%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateBottomBar(tt.width, tt.left, tt.center, tt.right)

			// Should not be empty
			if result == "" {
				t.Error("CreateBottomBar() returned empty string")
			}

			// Should contain non-empty sections
			if tt.left != "" && !strings.Contains(result, tt.left) {
				t.Errorf("CreateBottomBar() should contain left section %q", tt.left)
			}
			if tt.center != "" && !strings.Contains(result, tt.center) {
				t.Errorf("CreateBottomBar() should contain center section %q", tt.center)
			}
			if tt.right != "" && !strings.Contains(result, tt.right) {
				t.Errorf("CreateBottomBar() should contain right section %q", tt.right)
			}
		})
	}
}
