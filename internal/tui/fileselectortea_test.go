package tui

import (
	"flashcards/internal/tui/keys"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNewFileSelectorModel(t *testing.T) {
	files := []string{"/path/to/file1.md", "/path/to/file2.md"}
	model := NewFileSelectorModel(files)

	if model == nil {
		t.Fatal("NewFileSelectorModel() returned nil")
	}

	// Should have "All Files" option plus the provided files
	expectedLen := len(files) + 1
	if len(model.files) != expectedLen {
		t.Errorf("NewFileSelectorModel() files length = %d, want %d", len(model.files), expectedLen)
	}

	// First item should be "All Files"
	if model.files[0] != allFilesOption {
		t.Errorf("NewFileSelectorModel() first file = %q, want %q", model.files[0], allFilesOption)
	}

	// Cursor should start at 0
	if model.cursor != 0 {
		t.Errorf("NewFileSelectorModel() cursor = %d, want 0", model.cursor)
	}

	// No files should be selected initially
	if len(model.selected) != 0 {
		t.Errorf("NewFileSelectorModel() selected count = %d, want 0", len(model.selected))
	}
}

func TestFileSelectorModel_Init(t *testing.T) {
	model := NewFileSelectorModel([]string{"/path/to/file.md"})
	cmd := model.Init()

	if cmd != nil {
		t.Errorf("Init() = %v, want nil", cmd)
	}
}

func TestFileSelectorModel_Update_WindowSize(t *testing.T) {
	model := NewFileSelectorModel([]string{"/path/to/file.md"})

	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, _ := model.Update(msg)

	m := updatedModel.(*FileSelectorModel)
	if m.width != 100 {
		t.Errorf("Update(WindowSizeMsg) width = %d, want 100", m.width)
	}
	if m.height != 50 {
		t.Errorf("Update(WindowSizeMsg) height = %d, want 50", m.height)
	}
}

func TestFileSelectorModel_Update_Navigation(t *testing.T) {
	files := []string{"/file1.md", "/file2.md", "/file3.md"}

	tests := []struct {
		name           string
		key            string
		initialCursor  int
		expectedCursor int
	}{
		{
			name:           "down arrow",
			key:            keys.Down,
			initialCursor:  0,
			expectedCursor: 1,
		},
		{
			name:           "j key",
			key:            keys.J,
			initialCursor:  0,
			expectedCursor: 1,
		},
		{
			name:           "up arrow",
			key:            keys.Up,
			initialCursor:  1,
			expectedCursor: 0,
		},
		{
			name:           "k key",
			key:            keys.K,
			initialCursor:  1,
			expectedCursor: 0,
		},
		{
			name:           "up at top stays at top",
			key:            keys.Up,
			initialCursor:  0,
			expectedCursor: 0,
		},
		{
			name:           "down at bottom stays at bottom",
			key:            keys.Down,
			initialCursor:  3, // 3 files + "All Files" = 4 total, index 3 is last
			expectedCursor: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewFileSelectorModel(files)
			m.cursor = tt.initialCursor

			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{}, Alt: false}
			switch tt.key {
			case keys.Up:
				msg.Type = tea.KeyUp
			case keys.Down:
				msg.Type = tea.KeyDown
			default:
				msg.Runes = []rune(tt.key)
			}

			updatedModel, _ := m.Update(msg)
			updated := updatedModel.(*FileSelectorModel)

			if updated.cursor != tt.expectedCursor {
				t.Errorf("Update(%s) cursor = %d, want %d", tt.key, updated.cursor, tt.expectedCursor)
			}
		})
	}
}

func TestFileSelectorModel_Update_ToggleSelection(t *testing.T) {
	files := []string{"/file1.md", "/file2.md"}
	model := NewFileSelectorModel(files)
	model.cursor = 1 // Position on first actual file

	// Toggle selection with space
	msg := tea.KeyMsg{Type: tea.KeySpace}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(*FileSelectorModel)

	// First file should now be selected
	if !m.selected["/file1.md"] {
		t.Error("Update(Space) should select the file at cursor")
	}

	// Toggle again to deselect
	updatedModel2, _ := m.Update(msg)
	m2 := updatedModel2.(*FileSelectorModel)

	if m2.selected["/file1.md"] {
		t.Error("Update(Space) should deselect the file at cursor")
	}
}

func TestFileSelectorModel_Update_ToggleAllFiles(t *testing.T) {
	files := []string{"/file1.md", "/file2.md"}
	model := NewFileSelectorModel(files)
	model.cursor = 0 // Position on "All Files" option

	// Toggle "All Files" to select all
	msg := tea.KeyMsg{Type: tea.KeySpace}
	updatedModel, _ := model.Update(msg)
	m := updatedModel.(*FileSelectorModel)

	// Both files should be selected
	if !m.selected["/file1.md"] || !m.selected["/file2.md"] {
		t.Error("Update(Space on All Files) should select all files")
	}

	// Toggle again to deselect all
	updatedModel2, _ := m.Update(msg)
	m2 := updatedModel2.(*FileSelectorModel)

	if m2.selected["/file1.md"] || m2.selected["/file2.md"] {
		t.Error("Update(Space on All Files) should deselect all files")
	}
}

func TestFileSelectorModel_Update_Confirm(t *testing.T) {
	files := []string{"/file1.md", "/file2.md"}
	model := NewFileSelectorModel(files)
	model.selected["/file1.md"] = true

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, cmd := model.Update(msg)
	m := updatedModel.(*FileSelectorModel)

	if !m.confirmed {
		t.Error("Update(Enter) should set confirmed to true")
	}

	if len(m.selectedFiles) != 1 {
		t.Errorf("Update(Enter) selectedFiles length = %d, want 1", len(m.selectedFiles))
	}

	if cmd == nil {
		t.Error("Update(Enter) should return tea.Quit command")
	}
}

func TestFileSelectorModel_Update_Quit(t *testing.T) {
	tests := []struct {
		name string
		key  tea.KeyType
	}{
		{"quit with q", tea.KeyRunes},
		{"quit with ctrl+c", tea.KeyCtrlC},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewFileSelectorModel([]string{"/file1.md"})

			var msg tea.Msg
			if tt.key == tea.KeyRunes {
				msg = tea.KeyMsg{Type: tt.key, Runes: []rune("q")}
			} else {
				msg = tea.KeyMsg{Type: tt.key}
			}

			updatedModel, cmd := m.Update(msg)
			updated := updatedModel.(*FileSelectorModel)

			if !updated.confirmed {
				t.Error("Update(Quit) should set confirmed to true")
			}

			if len(updated.selectedFiles) != 0 {
				t.Errorf("Update(Quit) selectedFiles length = %d, want 0", len(updated.selectedFiles))
			}

			if cmd == nil {
				t.Error("Update(Quit) should return tea.Quit command")
			}
		})
	}
}

func TestFileSelectorModel_View(t *testing.T) {
	files := []string{"/file1.md", "/file2.md"}
	model := NewFileSelectorModel(files)
	model.width = 80
	model.height = 24

	view := model.View()

	if view == "" {
		t.Error("View() returned empty string")
	}

	// Should contain the "All Files" option
	if !contains(view, "All Files") {
		t.Error("View() should contain 'All Files' option")
	}
}

func TestFileSelectorModel_View_AfterConfirm(t *testing.T) {
	model := NewFileSelectorModel([]string{"/file1.md"})
	model.confirmed = true

	view := model.View()

	if view != "" {
		t.Errorf("View() after confirm = %q, want empty string", view)
	}
}

func TestFileSelectorModel_areAllFilesSelected(t *testing.T) {
	files := []string{"/file1.md", "/file2.md"}
	model := NewFileSelectorModel(files)

	// Initially, no files selected
	if model.areAllFilesSelected() {
		t.Error("areAllFilesSelected() should return false when no files selected")
	}

	// Select one file
	model.selected["/file1.md"] = true
	if model.areAllFilesSelected() {
		t.Error("areAllFilesSelected() should return false when not all files selected")
	}

	// Select all files
	model.selected["/file2.md"] = true
	if !model.areAllFilesSelected() {
		t.Error("areAllFilesSelected() should return true when all files selected")
	}
}

func TestFileSelectorModel_getSelectedFiles(t *testing.T) {
	files := []string{"/file1.md", "/file2.md", "/file3.md"}
	model := NewFileSelectorModel(files)

	// No files selected
	selected := model.getSelectedFiles()
	if len(selected) != 0 {
		t.Errorf("getSelectedFiles() length = %d, want 0", len(selected))
	}

	// Select some files
	model.selected["/file1.md"] = true
	model.selected["/file3.md"] = true

	selected = model.getSelectedFiles()
	if len(selected) != 2 {
		t.Errorf("getSelectedFiles() length = %d, want 2", len(selected))
	}

	// Should not include "All Files" option
	for _, f := range selected {
		if f == allFilesOption {
			t.Error("getSelectedFiles() should not include 'All Files' option")
		}
	}
}

func TestFileSelectorModel_GetSelectedFiles(t *testing.T) {
	files := []string{"/file1.md"}
	model := NewFileSelectorModel(files)
	model.selectedFiles = []string{"/file1.md"}

	result := model.GetSelectedFiles()

	if len(result) != 1 {
		t.Errorf("GetSelectedFiles() length = %d, want 1", len(result))
	}

	if result[0] != "/file1.md" {
		t.Errorf("GetSelectedFiles()[0] = %q, want %q", result[0], "/file1.md")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) >= len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
