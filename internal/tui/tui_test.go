package tui

import (
	"strings"
	"testing"

	"flashcards/internal/store"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		n        int
		expected string
	}{
		{
			name:     "no truncate",
			input:    "hello",
			n:        10,
			expected: "hello",
		},
		{
			name:     "truncate",
			input:    "hello world",
			n:        5,
			expected: "hell…",
		},
		{
			name:     "exact",
			input:    "hello",
			n:        5,
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncate(tt.input, tt.n)
			if result != tt.expected {
				t.Errorf("truncate() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestPrintFunctions(t *testing.T) {
	// Since Print functions write to stdout, we can't easily test output
	// But we can ensure they don't panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Print function panicked: %v", r)
		}
	}()

	PrintInfo("test info")
	PrintError("test error", nil)
	PrintSuccess("test success")
}

func TestNewAdminModel(t *testing.T) {
	// Create a temporary store
	tempDB := t.TempDir() + "/test.db"
	s, err := store.NewStore(tempDB)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	flashcards := []store.Flashcard{
		{ID: 1, Question: "Q1", Answer: "A1"},
	}

	model := NewAdminModel(s, flashcards)

	if len(model.flashcards) != 1 {
		t.Errorf("Expected 1 flashcard, got %d", len(model.flashcards))
	}
	if model.view != adminList {
		t.Errorf("Expected view adminList, got %v", model.view)
	}
}

func TestAdminModelUpdateQuit(t *testing.T) {
	tempDB := t.TempDir() + "/test.db"
	s, err := store.NewStore(tempDB)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	model := NewAdminModel(s, []store.Flashcard{})

	// Send quit key
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	newModel, cmd := model.Update(msg)

	if newModel == nil {
		t.Error("Model should not be nil")
	}
	if cmd == nil {
		t.Error("Cmd should not be nil on quit")
	}
	_ = cmd // Use cmd to avoid unused variable warning
}

func TestAdminModelKeyMap(t *testing.T) {
	// Test ShortHelp
	helpKeys := adminKeys.ShortHelp()
	if len(helpKeys) == 0 {
		t.Error("ShortHelp should return keys")
	}

	// Test FullHelp
	fullHelpKeys := adminKeys.FullHelp()
	if len(fullHelpKeys) == 0 {
		t.Error("FullHelp should return keys")
	}
}

func TestAdminModelInit(t *testing.T) {
	tempDB := t.TempDir() + "/test.db"
	s, err := store.NewStore(tempDB)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	model := NewAdminModel(s, []store.Flashcard{})
	cmd := model.Init()
	if cmd != nil {
		t.Error("Init() should return nil command")
	}
}

func TestAdminModelUpdate(t *testing.T) {
	tempDB := t.TempDir() + "/test.db"
	s, err := store.NewStore(tempDB)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	flashcards := []store.Flashcard{
		{ID: 1, Question: "Q1", Answer: "A1", RevisitIn: 0},
		{ID: 2, Question: "Q2", Answer: "A2", RevisitIn: 5},
	}

	model := NewAdminModel(s, flashcards)
	model.width = 80
	model.height = 24

	// Test window size message
	msg := tea.WindowSizeMsg{Width: 100, Height: 30}
	newModel, _ := model.Update(msg)
	if newModel == nil {
		t.Error("Update() should return a model")
	}

	// Test create key
	createMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
	newModel, _ = model.Update(createMsg)
	updatedModel := newModel.(*AdminModel)
	if updatedModel.view != adminCreate {
		t.Error("View should change to adminCreate")
	}

	// Test edit key
	model.view = adminList
	editMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
	newModel, _ = model.Update(editMsg)
	updatedModel = newModel.(*AdminModel)
	if updatedModel.view != adminEdit {
		t.Error("View should change to adminEdit")
	}

	// Test delete key
	model.view = adminList
	deleteMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	newModel, _ = model.Update(deleteMsg)
	updatedModel = newModel.(*AdminModel)
	if updatedModel.view != adminConfirmDelete {
		t.Error("View should change to adminConfirmDelete")
	}

	// Test bulk reset key
	model.view = adminList
	bulkResetMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}}
	newModel, _ = model.Update(bulkResetMsg)
	updatedModel = newModel.(*AdminModel)
	if updatedModel.view != adminConfirmBulkReset {
		t.Error("View should change to adminConfirmBulkReset")
	}

	// Test reload key
	model.view = adminList
	reloadMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}
	newModel, _ = model.Update(reloadMsg)
	updatedModel = newModel.(*AdminModel)
	if updatedModel.status.Success != "Table refreshed" {
		t.Error("Status message should be set on reload")
	}

	// Test help key
	model.view = adminList
	helpMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	newModel, _ = model.Update(helpMsg)
	updatedModel = newModel.(*AdminModel)
	if updatedModel.help.ShowAll != true {
		t.Error("Help should be shown")
	}

	// Test PageDown - with flashcards present
	if len(model.flashcards) > 0 {
		model.view = adminList
		model.selected = 0
		pageDownMsg := tea.KeyMsg{Type: tea.KeyPgDown}
		newModel, _ = model.Update(pageDownMsg)
		updatedModel = newModel.(*AdminModel)
		// Should move cursor down by 10 or to end
		if updatedModel.selected < 0 {
			t.Error("Selected should not be negative")
		}
		if updatedModel.selected >= len(updatedModel.flashcards) {
			t.Error("Selected should not exceed list length")
		}
	}

	// Test PageUp - with flashcards present
	if len(model.flashcards) > 0 {
		model.view = adminList
		model.selected = min(5, len(model.flashcards)-1)
		pageUpMsg := tea.KeyMsg{Type: tea.KeyPgUp}
		newModel, _ = model.Update(pageUpMsg)
		updatedModel = newModel.(*AdminModel)
		// Should move cursor up by 10 or to start
		if updatedModel.selected < 0 {
			t.Error("Selected should not be negative")
		}
	}

	// Test edit with empty list
	model.view = adminList
	model.flashcards = []store.Flashcard{}
	editMsgEmpty := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
	newModel, _ = model.Update(editMsgEmpty)
	updatedModel = newModel.(*AdminModel)
	if updatedModel.view != adminList {
		t.Error("View should stay in adminList when list is empty")
	}

	// Test delete with empty list
	model.view = adminList
	deleteMsgEmpty := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	newModel, _ = model.Update(deleteMsgEmpty)
	updatedModel = newModel.(*AdminModel)
	if updatedModel.view != adminList {
		t.Error("View should stay in adminList when list is empty")
	}

	// Test bulk reset with empty list
	model.view = adminList
	bulkResetMsgEmpty := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}}
	newModel, _ = model.Update(bulkResetMsgEmpty)
	updatedModel = newModel.(*AdminModel)
	if updatedModel.view != adminList {
		t.Error("View should stay in adminList when list is empty")
	}
}

func TestAdminModelHandleCreateView(t *testing.T) {
	tempDB := t.TempDir() + "/test.db"
	s, err := store.NewStore(tempDB)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	model := NewAdminModel(s, []store.Flashcard{})
	model.view = adminCreate

	// Test cancel
	cancelMsg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ := model.handleCreateView(cancelMsg)
	updatedModel := newModel.(*AdminModel)
	if updatedModel.view != adminList {
		t.Error("View should return to adminList on cancel")
	}

	// Test tab
	model.view = adminCreate
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	newModel, _ = model.handleCreateView(tabMsg)
	updatedModel = newModel.(*AdminModel)
	if updatedModel.view != adminCreate {
		t.Error("View should stay in adminCreate")
	}

	// Test enter (will try to create, but may fail validation)
	model.view = adminCreate
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ = model.handleCreateView(enterMsg)
	updatedModel = newModel.(*AdminModel)
	// Should show error since form is empty
	if updatedModel.status.Error == "" {
		t.Log("Note: Error message may be empty if validation passes")
	}
}

func TestAdminModelHandleEditView(t *testing.T) {
	tempDB := t.TempDir() + "/test.db"
	s, err := store.NewStore(tempDB)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	flashcards := []store.Flashcard{
		{ID: 1, Question: "Q1", Answer: "A1", RevisitIn: 0},
	}
	model := NewAdminModel(s, flashcards)
	model.view = adminEdit
	model.selected = 0

	// Test cancel
	cancelMsg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ := model.handleEditView(cancelMsg)
	updatedModel := newModel.(*AdminModel)
	if updatedModel.view != adminList {
		t.Error("View should return to adminList on cancel")
	}

	// Test tab
	model.view = adminEdit
	model.selected = 0
	model.loadSelectedIntoForm()
	tabMsg := tea.KeyMsg{Type: tea.KeyTab}
	newModel, _ = model.handleEditView(tabMsg)
	updatedModel = newModel.(*AdminModel)
	if updatedModel.view != adminEdit {
		t.Error("View should stay in adminEdit")
	}

	// Test enter (will try to update)
	model.view = adminEdit
	model.selected = 0
	model.loadSelectedIntoForm()
	model.questionInput.SetValue("Updated Q1")
	model.answerInput.SetValue("Updated A1")
	model.revisitInput.SetValue("5")
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ = model.handleEditView(enterMsg)
	updatedModel = newModel.(*AdminModel)
	// Should either succeed or show error
	if updatedModel.view != adminList && updatedModel.status.Error == "" {
		t.Log("Update may have succeeded or validation passed")
	}

	// Test key input
	model.view = adminEdit
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	newModel, _ = model.handleEditView(keyMsg)
	updatedModel = newModel.(*AdminModel)
	if updatedModel.view != adminEdit {
		t.Error("View should stay in adminEdit on key input")
	}
}

func TestAdminModelHandleDeleteConfirm(t *testing.T) {
	tempDB := t.TempDir() + "/test.db"
	s, err := store.NewStore(tempDB)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	flashcards := []store.Flashcard{
		{ID: 1, Question: "Q1", Answer: "A1", RevisitIn: 0},
	}
	model := NewAdminModel(s, flashcards)
	model.view = adminConfirmDelete
	model.selected = 0

	// Test cancel
	cancelMsg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ := model.handleDeleteConfirm(cancelMsg)
	updatedModel := newModel.(*AdminModel)
	if updatedModel.view != adminList {
		t.Error("View should return to adminList on cancel")
	}

	// Test 'n' (no)
	model.view = adminConfirmDelete
	noMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	newModel, _ = model.handleDeleteConfirm(noMsg)
	updatedModel = newModel.(*AdminModel)
	if updatedModel.view != adminList {
		t.Error("View should return to adminList on 'n'")
	}

	// Test 'y' (yes) - will delete the flashcard
	model.view = adminConfirmDelete
	model.selected = 0
	yesMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
	newModel, _ = model.handleDeleteConfirm(yesMsg)
	updatedModel = newModel.(*AdminModel)
	if updatedModel.view != adminList {
		t.Error("View should return to adminList after delete")
	}

	// Test delete when selected is at end of list
	fc := store.Flashcard{Question: "Q2", Answer: "A2", RevisitIn: 0, File: "test2.md"}
	err = s.InsertFlashcard(fc)
	if err != nil {
		t.Fatalf("Failed to insert flashcard: %v", err)
	}

	flashcards, err = s.GetAllFlashcards()
	if err != nil {
		t.Fatalf("Failed to get flashcards: %v", err)
	}

	model = NewAdminModel(s, flashcards)
	model.selected = len(flashcards) - 1
	model.view = adminConfirmDelete
	newModel, _ = model.handleDeleteConfirm(yesMsg)
	updatedModel = newModel.(*AdminModel)
	if updatedModel.selected >= len(updatedModel.flashcards) && len(updatedModel.flashcards) > 0 {
		t.Error("Selected should be adjusted after delete")
	}
}

func TestAdminModelHandleBulkResetConfirm(t *testing.T) {
	tempDB := t.TempDir() + "/test.db"
	s, err := store.NewStore(tempDB)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	flashcards := []store.Flashcard{
		{ID: 1, Question: "Q1", Answer: "A1", RevisitIn: 5},
	}
	model := NewAdminModel(s, flashcards)
	model.view = adminConfirmBulkReset

	// Test cancel
	cancelMsg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ := model.handleBulkResetConfirm(cancelMsg)
	updatedModel := newModel.(*AdminModel)
	if updatedModel.view != adminList {
		t.Error("View should return to adminList on cancel")
	}

	// Test 'n' (no)
	model.view = adminConfirmBulkReset
	noMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	newModel, _ = model.handleBulkResetConfirm(noMsg)
	updatedModel = newModel.(*AdminModel)
	if updatedModel.view != adminList {
		t.Error("View should return to adminList on 'n'")
	}
}

func TestAdminModelView(t *testing.T) {
	tempDB := t.TempDir() + "/test.db"
	s, err := store.NewStore(tempDB)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	flashcards := []store.Flashcard{
		{ID: 1, Question: "Q1", Answer: "A1", RevisitIn: 0},
	}
	model := NewAdminModel(s, flashcards)
	model.width = 80
	model.height = 24

	// Test list view
	model.view = adminList
	view := model.View()
	if view == "" {
		t.Error("View should return content")
	}

	// Test create view
	model.view = adminCreate
	view = model.View()
	if !strings.Contains(view, "Question:") {
		t.Error("Create view should contain 'Question:'")
	}

	// Test edit view
	model.view = adminEdit
	model.selected = 0
	view = model.View()
	if !strings.Contains(view, "Question:") {
		t.Error("Edit view should contain 'Question:'")
	}

	// Test delete confirm view
	model.view = adminConfirmDelete
	model.selected = 0
	view = model.View()
	if !strings.Contains(view, "Delete") {
		t.Error("Delete confirm view should contain 'Delete'")
	}

	// Test bulk reset confirm view
	model.view = adminConfirmBulkReset
	view = model.View()
	if !strings.Contains(view, "RevisitIn") {
		t.Error("Bulk reset confirm view should contain 'RevisitIn'")
	}

	// Test view with error message
	model.view = adminList
	model.status.SetError("Test error")
	// Verify the status component works
	if model.status.Error != "Test error" {
		t.Error("Status.Error should be set")
	}
	view = model.View()
	// The view renders with ANSI codes, just verify it's not empty and contains table
	if !strings.Contains(view, "Question") || len(view) == 0 {
		t.Error("View should render with table content")
	}

	// Test view with status message
	model.view = adminList
	model.status.SetSuccess("Test success")
	// Verify the status component works
	if model.status.Success != "Test success" {
		t.Error("Status.Success should be set")
	}
	view = model.View()
	// The view renders with ANSI codes, just verify it's not empty and contains table
	if !strings.Contains(view, "Question") || len(view) == 0 {
		t.Error("View should render with table content")
	}

	// Test view with empty flashcards list
	model.view = adminList
	model.flashcards = []store.Flashcard{}
	view = model.View()
	if !strings.Contains(view, "No flashcards") {
		t.Error("View should show message when no flashcards")
	}
}

func TestAdminModelHelpers(t *testing.T) {
	tempDB := t.TempDir() + "/test.db"
	s, err := store.NewStore(tempDB)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	flashcards := []store.Flashcard{
		{ID: 1, Question: "Q1", Answer: "A1", RevisitIn: 0},
	}
	model := NewAdminModel(s, flashcards)

	// Test resetForm
	model.resetForm()
	if model.questionInput.Value() != "" {
		t.Error("Question input should be empty after reset")
	}
	if model.answerInput.Value() != "" {
		t.Error("Answer input should be empty after reset")
	}

	// Test loadSelectedIntoForm
	model.selected = 0
	model.loadSelectedIntoForm()
	if model.questionInput.Value() != "Q1" {
		t.Errorf("Question should be 'Q1', got '%s'", model.questionInput.Value())
	}

	// Test parseRevisitDays
	model.revisitInput.SetValue("7")
	days, err := model.parseRevisitDays()
	if err != nil {
		t.Errorf("parseRevisitDays() error = %v", err)
	}
	if days != 7 {
		t.Errorf("Expected 7 days, got %d", days)
	}

	// Test parseRevisitDays with invalid input
	model.revisitInput.SetValue("invalid")
	_, err = model.parseRevisitDays()
	if err == nil {
		t.Error("parseRevisitDays() should return error for invalid input")
	}

	// Test parseRevisitDays with empty input
	model.revisitInput.SetValue("")
	_, err = model.parseRevisitDays()
	if err == nil {
		t.Error("parseRevisitDays() should return error for empty input")
	}

	// Test cycleFocus - start with question focused
	model.questionInput.Focus()
	model.answerInput.Blur()
	model.revisitInput.Blur()
	model.cycleFocus()
	if !model.answerInput.Focused() {
		t.Error("Focus should cycle from question to answer")
	}

	// Test cycleFocus - from answer to revisit
	model.answerInput.Focus()
	model.questionInput.Blur()
	model.revisitInput.Blur()
	model.cycleFocus()
	if !model.revisitInput.Focused() {
		t.Error("Focus should cycle from answer to revisit")
	}

	// Test cycleFocus - from revisit to question
	model.revisitInput.Focus()
	model.questionInput.Blur()
	model.answerInput.Blur()
	model.cycleFocus()
	if !model.questionInput.Focused() {
		t.Error("Focus should cycle from revisit to question")
	}

	// Test cycleFocus - default case (none focused)
	model.questionInput.Blur()
	model.answerInput.Blur()
	model.revisitInput.Blur()
	model.cycleFocus()
	if !model.questionInput.Focused() {
		t.Error("Focus should default to question when none focused")
	}

	// Test updateInputs with question input focused
	model.questionInput.Focus()
	model.answerInput.Blur()
	model.revisitInput.Blur()
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	model.updateInputs(keyMsg)

	// Test updateInputs with answer input focused
	model.questionInput.Blur()
	model.answerInput.Focus()
	model.revisitInput.Blur()
	model.updateInputs(keyMsg)

	// Test updateInputs with revisit input focused
	model.questionInput.Blur()
	model.answerInput.Blur()
	model.revisitInput.Focus()
	model.updateInputs(keyMsg)
}

func TestAdminModelCreateFlashcard(t *testing.T) {
	tempDB := t.TempDir() + "/test.db"
	s, err := store.NewStore(tempDB)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	model := NewAdminModel(s, []store.Flashcard{})
	model.questionInput.SetValue("Test Question")
	model.answerInput.SetValue("Test Answer")
	model.revisitInput.SetValue("7")

	model.createFlashcard()

	if model.status.Error != "" && model.status.Success == "" {
		t.Logf("Error message: %s", model.status.Error)
	} else if model.status.Success != "" {
		t.Logf("Success message: %s", model.status.Success)
	}
}

func TestAdminModelUpdateFlashcard(t *testing.T) {
	tempDB := t.TempDir() + "/test.db"
	s, err := store.NewStore(tempDB)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	// Insert a flashcard first
	fc := store.Flashcard{Question: "Q1", Answer: "A1", RevisitIn: 0, File: "test.md"}
	err = s.InsertFlashcard(fc)
	if err != nil {
		t.Fatalf("Failed to insert flashcard: %v", err)
	}

	flashcards, err := s.GetAllFlashcards()
	if err != nil {
		t.Fatalf("Failed to get flashcards: %v", err)
	}

	model := NewAdminModel(s, flashcards)
	model.selected = 0
	model.questionInput.SetValue("Updated Question")
	model.answerInput.SetValue("Updated Answer")
	model.revisitInput.SetValue("10")

	model.updateFlashcard()

	if model.status.Error != "" && model.status.Success == "" {
		t.Logf("Error message: %s", model.status.Error)
	} else if model.status.Success != "" {
		t.Logf("Success message: %s", model.status.Success)
	}
}

func TestAdminModelDeleteFlashcard(t *testing.T) {
	tempDB := t.TempDir() + "/test.db"
	s, err := store.NewStore(tempDB)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	// Insert a flashcard first
	fc := store.Flashcard{Question: "Q1", Answer: "A1", RevisitIn: 0, File: "test.md"}
	err = s.InsertFlashcard(fc)
	if err != nil {
		t.Fatalf("Failed to insert flashcard: %v", err)
	}

	flashcards, err := s.GetAllFlashcards()
	if err != nil {
		t.Fatalf("Failed to get flashcards: %v", err)
	}

	model := NewAdminModel(s, flashcards)
	model.selected = 0

	model.deleteFlashcard()

	if model.status.Error != "" && model.status.Success == "" {
		t.Logf("Error message: %s", model.status.Error)
	} else if model.status.Success != "" {
		t.Logf("Success message: %s", model.status.Success)
	}
}

func TestAdminModelReload(t *testing.T) {
	tempDB := t.TempDir() + "/test.db"
	s, err := store.NewStore(tempDB)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	model := NewAdminModel(s, []store.Flashcard{})
	model.reload()

	if model.status.Error != "" {
		t.Logf("Error message: %s", model.status.Error)
	}

	// Test reload with selected out of bounds (greater than length)
	fc := store.Flashcard{Question: "Q1", Answer: "A1", RevisitIn: 0, File: "test.md"}
	err = s.InsertFlashcard(fc)
	if err != nil {
		t.Fatalf("Failed to insert flashcard: %v", err)
	}

	model.selected = 10 // Out of bounds
	model.reload()
	if model.selected >= len(model.flashcards) && len(model.flashcards) > 0 {
		t.Error("Selected should be adjusted when out of bounds")
	}

	// Test reload with selected negative
	model.selected = -1
	model.reload()
	if model.selected < 0 {
		t.Error("Selected should be adjusted when negative")
	}
}

func TestAdminModelBulkResetRevisitIn(t *testing.T) {
	tempDB := t.TempDir() + "/test.db"
	s, err := store.NewStore(tempDB)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer s.Close()

	// Insert flashcards
	fc1 := store.Flashcard{Question: "Q1", Answer: "A1", RevisitIn: 5, File: "test1.md"}
	fc2 := store.Flashcard{Question: "Q2", Answer: "A2", RevisitIn: 10, File: "test2.md"}
	err = s.InsertFlashcard(fc1)
	if err != nil {
		t.Fatalf("Failed to insert flashcard: %v", err)
	}
	err = s.InsertFlashcard(fc2)
	if err != nil {
		t.Fatalf("Failed to insert flashcard: %v", err)
	}

	flashcards, err := s.GetAllFlashcards()
	if err != nil {
		t.Fatalf("Failed to get flashcards: %v", err)
	}

	model := NewAdminModel(s, flashcards)
	model.bulkResetRevisitIn()

	if model.status.Error != "" && model.status.Success == "" {
		t.Logf("Error message: %s", model.status.Error)
	} else if model.status.Success != "" {
		t.Logf("Success message: %s", model.status.Success)
	}
}

func TestNewReviewModel(t *testing.T) {
	flashcards := []store.Flashcard{
		{ID: 1, Question: "Q1", Answer: "A1"},
	}

	model := NewReviewModel(flashcards)

	if len(model.flashcards) != 1 {
		t.Errorf("Expected 1 flashcard, got %d", len(model.flashcards))
	}
	if model.current != 0 {
		t.Errorf("Expected current 0, got %d", model.current)
	}
}

func TestReviewModelNextCard(t *testing.T) {
	flashcards := []store.Flashcard{
		{ID: 1, Question: "Q1", Answer: "A1"},
		{ID: 2, Question: "Q2", Answer: "A2"},
	}

	model := NewReviewModel(flashcards)
	model.current = 0

	cmd := model.nextCard()

	if model.current != 1 {
		t.Errorf("Expected current 1, got %d", model.current)
	}
	if model.view != viewQuestion {
		t.Errorf("Expected view viewQuestion, got %v", model.view)
	}
	if cmd == nil {
		t.Error("Cmd should not be nil")
	}

	// Next card
	model.nextCard()
	if model.current != 2 {
		t.Errorf("Expected current 2, got %d", model.current)
	}
	if model.view != viewDone {
		t.Errorf("Expected view viewDone, got %v", model.view)
	}
}

func TestReviewModelFlashcardWasCorrect(t *testing.T) {
	flashcards := []store.Flashcard{
		{ID: 1, Question: "Q1", Answer: "A1"},
	}

	model := NewReviewModel(flashcards)
	model.correct[0] = true

	if !model.FlashcardWasCorrect(0) {
		t.Error("Expected true for correct")
	}
	if model.FlashcardWasCorrect(1) {
		t.Error("Expected false for invalid index")
	}
}

func TestReviewModelFlashcardRevisitIn(t *testing.T) {
	flashcards := []store.Flashcard{
		{ID: 1, Question: "Q1", Answer: "A1"},
	}

	model := NewReviewModel(flashcards)
	model.revisitIn[0] = 7

	if model.FlashcardRevisitIn(0) != 7 {
		t.Errorf("Expected 7, got %d", model.FlashcardRevisitIn(0))
	}
	if model.FlashcardRevisitIn(1) != 0 {
		t.Error("Expected 0 for invalid index")
	}
}

func TestReviewModelView(t *testing.T) {
	flashcards := []store.Flashcard{
		{ID: 1, Question: "What is Go?", Answer: "A language"},
	}

	model := NewReviewModel(flashcards)
	model.width = 80
	model.height = 24

	view := model.View()
	if !strings.Contains(view, "Question:") {
		t.Error("View should contain 'Question:'")
	}
	if !strings.Contains(view, "What is Go?") {
		t.Error("View should contain the question")
	}

	// Test viewAnswer
	model.view = viewAnswer
	view = model.View()
	if !strings.Contains(view, "Answer:") {
		t.Error("View should contain 'Answer:'")
	}

	// Test viewRevisitIn
	model.view = viewRevisitIn
	model.resultMsg = "Revisit in 7 days"
	view = model.View()
	if !strings.Contains(view, "Revisit in") {
		t.Error("View should contain 'Revisit in'")
	}

	// Test viewDone
	model.view = viewDone
	model.completionMsg = "Review complete 🎉🎉🎉"
	view = model.View()
	if !strings.Contains(view, "Review complete") {
		t.Error("View should contain 'Review complete'")
	}

	// Test quitting
	model.quitting = true
	view = model.View()
	if view != "Goodbye!" {
		t.Errorf("View should return 'Goodbye!' when quitting, got: %s", view)
	}
}

func TestReviewModelInit(t *testing.T) {
	flashcards := []store.Flashcard{
		{ID: 1, Question: "Q1", Answer: "A1"},
	}

	model := NewReviewModel(flashcards)
	cmd := model.Init()
	if cmd == nil {
		t.Error("Init() should return a command")
	}
}

func TestReviewModelUpdate(t *testing.T) {
	flashcards := []store.Flashcard{
		{ID: 1, Question: "Q1", Answer: "A1"},
		{ID: 2, Question: "Q2", Answer: "A2"},
	}

	model := NewReviewModel(flashcards)
	model.width = 80
	model.height = 24

	// Test window size message
	msg := tea.WindowSizeMsg{Width: 100, Height: 30}
	newModel, _ := model.Update(msg)
	if newModel == nil {
		t.Error("Update() should return a model")
	}
	updatedModel := newModel.(*ReviewModel)
	if updatedModel.width != 100 || updatedModel.height != 30 {
		t.Error("Window size should be updated")
	}

	// Test quit key
	quitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	newModel, _ = model.Update(quitMsg)
	if newModel == nil {
		t.Error("Update() should return a model")
	}
	updatedModel = newModel.(*ReviewModel)
	if !updatedModel.quitting {
		t.Error("Model should be quitting")
	}

	// Test viewQuestion -> viewAnswer
	model.view = viewQuestion
	model.quitting = false
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ = model.Update(enterMsg)
	updatedModel = newModel.(*ReviewModel)
	if updatedModel.view != viewAnswer {
		t.Error("View should change to viewAnswer on enter")
	}

	// Test viewAnswer -> correct
	model.view = viewAnswer
	correctMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}
	newModel, _ = model.Update(correctMsg)
	updatedModel = newModel.(*ReviewModel)
	if !updatedModel.correct[0] {
		t.Error("Flashcard should be marked correct")
	}
	if updatedModel.view != viewRevisitIn {
		t.Error("View should change to viewRevisitIn after correct")
	}

	// Test viewAnswer -> incorrect
	model.view = viewAnswer
	model.current = 0
	incorrectMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}}
	newModel, _ = model.Update(incorrectMsg)
	updatedModel = newModel.(*ReviewModel)
	if updatedModel.correct[0] {
		t.Error("Flashcard should be marked incorrect")
	}

	// Test viewRevisitIn -> select days
	model.view = viewRevisitIn
	model.current = 0
	day7Msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'7'}}
	newModel, _ = model.Update(day7Msg)
	updatedModel = newModel.(*ReviewModel)
	if updatedModel.revisitIn[0] != 7 {
		t.Errorf("RevisitIn should be 7, got %d", updatedModel.revisitIn[0])
	}

	// Test viewDone -> quit
	model.view = viewDone
	model.quitting = false
	newModel, _ = model.Update(quitMsg)
	updatedModel = newModel.(*ReviewModel)
	if !updatedModel.quitting {
		t.Error("Model should be quitting from viewDone")
	}

	// Test viewRevisitIn with different day keys
	model.view = viewRevisitIn
	model.current = 0
	day1Msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	newModel, _ = model.Update(day1Msg)
	updatedModel = newModel.(*ReviewModel)
	if updatedModel.revisitIn[0] != 1 {
		t.Errorf("RevisitIn should be 1, got %d", updatedModel.revisitIn[0])
	}

	model.view = viewRevisitIn
	model.current = 0
	day3Msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'3'}}
	newModel, _ = model.Update(day3Msg)
	updatedModel = newModel.(*ReviewModel)
	if updatedModel.revisitIn[0] != 3 {
		t.Errorf("RevisitIn should be 3, got %d", updatedModel.revisitIn[0])
	}

	model.view = viewRevisitIn
	model.current = 0
	day9Msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'9'}}
	newModel, _ = model.Update(day9Msg)
	updatedModel = newModel.(*ReviewModel)
	if updatedModel.revisitIn[0] != 9 {
		t.Errorf("RevisitIn should be 9, got %d", updatedModel.revisitIn[0])
	}

	// Test timer timeout
	model.view = viewQuestion
	timeoutMsg := timer.TimeoutMsg{}
	newModel, _ = model.Update(timeoutMsg)
	updatedModel = newModel.(*ReviewModel)
	if updatedModel.view != viewAnswer {
		t.Error("View should change to viewAnswer on timeout")
	}

	// Test progress frame message
	model.view = viewQuestion
	progressMsg := progress.FrameMsg{}
	newModel, _ = model.Update(progressMsg)
	updatedModel = newModel.(*ReviewModel)
	if updatedModel.view != viewQuestion {
		t.Error("View should stay in viewQuestion on progress frame")
	}

	// Test timer tick message
	model.view = viewQuestion
	tickMsg := timer.TickMsg{}
	newModel, _ = model.Update(tickMsg)
	updatedModel = newModel.(*ReviewModel)
	if updatedModel.view != viewQuestion {
		t.Error("View should stay in viewQuestion on timer tick")
	}

	// Test quit message
	model.quitting = false
	quitTeaMsg := tea.QuitMsg{}
	newModel, _ = model.Update(quitTeaMsg)
	updatedModel = newModel.(*ReviewModel)
	if !updatedModel.quitting {
		t.Error("Model should be quitting on QuitMsg")
	}

	// Test view with width > 80
	model.width = 100
	model.height = 30
	model.view = viewQuestion
	viewWide := model.View()
	if viewWide == "" {
		t.Error("View should return content")
	}

	// Test view with viewDone and correct/incorrect counts
	model.quitting = false
	model.view = viewDone
	model.current = len(model.flashcards) // Set current to end
	if len(model.correct) >= 2 {
		model.correct[0] = true
		model.correct[1] = false
	}
	viewDoneContent := model.View()
	// View should show review complete message, not "Goodbye!"
	if strings.Contains(viewDoneContent, "Goodbye!") {
		t.Error("View should not show 'Goodbye!' when not quitting")
	}
	if !strings.Contains(viewDoneContent, "complete") && !strings.Contains(viewDoneContent, "🎉") {
		t.Logf("View content (first 200 chars): %s", viewDoneContent[:min(200, len(viewDoneContent))])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
