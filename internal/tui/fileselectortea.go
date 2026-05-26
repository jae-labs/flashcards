package tui

import (
	"flashcards/internal/tui/keys"
	"flashcards/internal/tui/layout"
	"flashcards/internal/tui/theme"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	allFilesOption = "📚 All Files"
)

// FileSelectorModel represents the file selection screen
type FileSelectorModel struct {
	files         []string        // List of file paths
	selected      map[string]bool // Map of selected files
	cursor        int             // Current cursor position
	width         int             // Terminal width
	height        int             // Terminal height
	confirmed     bool            // Whether user confirmed selection
	selectedFiles []string        // Final selected files after confirmation
}

// NewFileSelectorModel creates a new file selector model
func NewFileSelectorModel(files []string) *FileSelectorModel {
	selected := make(map[string]bool)
	// Add "All Files" option at the beginning
	allFiles := append([]string{allFilesOption}, files...)

	return &FileSelectorModel{
		files:    allFiles,
		selected: selected,
		cursor:   0,
	}
}

func (m *FileSelectorModel) Init() tea.Cmd {
	return nil
}

func (m *FileSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case keys.CtrlC, keys.Q:
			// User quit - return empty selection
			m.confirmed = true
			m.selectedFiles = []string{}
			return m, tea.Quit

		case keys.Up, keys.K:
			if m.cursor > 0 {
				m.cursor--
			}

		case keys.Down, keys.J:
			if m.cursor < len(m.files)-1 {
				m.cursor++
			}

		case keys.Space: // Spacebar to toggle selection
			currentFile := m.files[m.cursor]
			if currentFile == allFilesOption {
				// Toggle all files
				allSelected := m.areAllFilesSelected()
				if allSelected {
					// Deselect all
					m.selected = make(map[string]bool)
				} else {
					// Select all
					for _, f := range m.files {
						if f != allFilesOption {
							m.selected[f] = true
						}
					}
				}
			} else {
				// Toggle individual file
				m.selected[currentFile] = !m.selected[currentFile]
			}

		case keys.Enter:
			// Confirm selection
			m.confirmed = true
			m.selectedFiles = m.getSelectedFiles()
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *FileSelectorModel) View() string {
	if m.confirmed {
		return ""
	}

	// Use layout helper for content width
	width := layout.CalculateContentWidth(m.width)

	var s strings.Builder

	// Title - centered with padding
	titleWithPadding := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(theme.ColorPrimary)).
		Padding(1, 0)
	title := titleWithPadding.Render("Welcome 👋 please select file(s) to get started!")
	centeredTitle := lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(title)
	s.WriteString(centeredTitle)
	s.WriteString("\n")

	// Calculate visible area - need to account for title, help, and selected count
	maxVisible := 15 // Fixed reasonable number of visible items
	if m.height > 20 {
		maxVisible = m.height - 15 // Adjust based on screen size
	}
	if maxVisible < 5 {
		maxVisible = 5
	}

	// Calculate scroll offset
	scrollOffset := 0
	if m.cursor >= maxVisible {
		scrollOffset = m.cursor - maxVisible + 1
	}

	// Render file list
	for i := scrollOffset; i < len(m.files) && i < scrollOffset+maxVisible; i++ {
		file := m.files[i]
		cursor := " "
		if i == m.cursor {
			cursor = theme.CursorStyle.Render("❯")
		}

		checkbox := theme.UncheckedStyle.Render("☐")
		itemStyle := theme.UnselectedStyle

		if file == allFilesOption {
			// Check if all files are selected
			if m.areAllFilesSelected() {
				checkbox = theme.CheckedStyle.Render("☑")
				itemStyle = theme.SelectedStyle
			}
		} else {
			// Regular file
			if m.selected[file] {
				checkbox = theme.CheckedStyle.Render("☑")
				itemStyle = theme.SelectedStyle
			}
		}

		// Truncate long file paths for display
		displayFile := file
		maxWidth := width - 10
		if maxWidth < 20 {
			maxWidth = 20
		}
		if len(displayFile) > maxWidth {
			displayFile = "..." + displayFile[len(displayFile)-maxWidth+3:]
		}

		fmt.Fprintf(&s, "%s %s %s\n", cursor, checkbox, itemStyle.Render(displayFile))
	}

	// Show scroll indicator if needed
	if len(m.files) > maxVisible {
		s.WriteString(theme.InfoStyle.Render(fmt.Sprintf("\n(Showing %d-%d of %d files)",
			scrollOffset+1,
			min(scrollOffset+maxVisible, len(m.files)),
			len(m.files))))
		s.WriteString("\n")
	}

	// Count selected files
	selectedCount := len(m.getSelectedFiles())
	if selectedCount > 0 {
		s.WriteString("\n")
		s.WriteString(theme.SuccessStyle.Render(fmt.Sprintf("Selected: %d file(s)", selectedCount)))
	}

	// Create frame with border using layout helper
	frame := layout.CreateFrame(width,
		layout.WithAlignment(lipgloss.Left, lipgloss.Top),
		layout.WithPadding(1, 2))

	// Help text (outside the frame)
	exitMsg := theme.InfoStyle.Render("↑/↓: Navigate • Space: Toggle • Enter: Confirm • q: Quit")

	// Center everything on screen using layout helper
	return layout.CenterContent(m.width, m.height, frame.Render(s.String())+"\n"+exitMsg)
}

// areAllFilesSelected checks if all files (excluding "All Files" option) are selected
func (m *FileSelectorModel) areAllFilesSelected() bool {
	if len(m.selected) == 0 {
		return false
	}
	for _, f := range m.files {
		if f != allFilesOption && !m.selected[f] {
			return false
		}
	}
	return true
}

// getSelectedFiles returns the list of selected file paths
func (m *FileSelectorModel) getSelectedFiles() []string {
	var selected []string
	for _, f := range m.files {
		if f != allFilesOption && m.selected[f] {
			selected = append(selected, f)
		}
	}
	return selected
}

// GetSelectedFiles returns the final selected files after confirmation
func (m *FileSelectorModel) GetSelectedFiles() []string {
	return m.selectedFiles
}
