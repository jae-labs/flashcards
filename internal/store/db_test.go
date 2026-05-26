package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewStore(t *testing.T) {
	// Create a temporary database file
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("NewStore() error = %v", err)
	}
	defer store.Close()

	// Check if database file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}
}

func TestInsertFlashcard(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	flashcard := Flashcard{
		File:      "/test/file.md",
		Question:  "What is 2+2?",
		Answer:    "4",
		RevisitIn: 0,
	}

	err := store.InsertFlashcard(flashcard)
	if err != nil {
		t.Errorf("InsertFlashcard() error = %v", err)
	}

	// Verify the flashcard was inserted
	cards, err := store.GetAllFlashcards()
	if err != nil {
		t.Fatalf("GetAllFlashcards() error = %v", err)
	}

	if len(cards) != 1 {
		t.Errorf("Expected 1 flashcard, got %d", len(cards))
	}

	if cards[0].Question != flashcard.Question {
		t.Errorf("Expected question '%s', got '%s'", flashcard.Question, cards[0].Question)
	}
}

func TestGetFlashcardsForReview(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	// Insert test flashcards
	flashcards := []Flashcard{
		{File: "/test/1.md", Question: "Q1", Answer: "A1", RevisitIn: -1}, // Due
		{File: "/test/2.md", Question: "Q2", Answer: "A2", RevisitIn: 5},  // Not due
		{File: "/test/3.md", Question: "Q3", Answer: "A3", RevisitIn: 0},  // Due
	}

	for _, fc := range flashcards {
		err := store.InsertFlashcard(fc)
		if err != nil {
			t.Fatalf("InsertFlashcard() error = %v", err)
		}
	}

	// Get flashcards for review
	dueCards, err := store.GetFlashcardsForReview()
	if err != nil {
		t.Fatalf("GetFlashcardsForReview() error = %v", err)
	}

	if len(dueCards) != 2 {
		t.Errorf("Expected 2 flashcards due for review, got %d", len(dueCards))
	}
}

func TestUpdateFlashcard(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	// Insert a flashcard
	fc := Flashcard{File: "/test/1.md", Question: "Q1", Answer: "A1", RevisitIn: 0}
	err := store.InsertFlashcard(fc)
	if err != nil {
		t.Fatalf("InsertFlashcard() error = %v", err)
	}

	// Get the inserted flashcard
	cards, err := store.GetAllFlashcards()
	if err != nil {
		t.Fatalf("GetAllFlashcards() error = %v", err)
	}

	if len(cards) != 1 {
		t.Fatalf("Expected 1 flashcard, got %d", len(cards))
	}

	// Update the flashcard
	cards[0].RevisitIn = 5
	err = store.UpdateFlashcard(cards[0])
	if err != nil {
		t.Errorf("UpdateFlashcard() error = %v", err)
	}

	// Verify the update
	updatedCards, err := store.GetAllFlashcards()
	if err != nil {
		t.Fatalf("GetAllFlashcards() error = %v", err)
	}

	if updatedCards[0].RevisitIn != 5 {
		t.Errorf("Expected RevisitIn 5, got %d", updatedCards[0].RevisitIn)
	}
}

func TestDeleteFlashcard(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	// Insert a flashcard
	fc := Flashcard{File: "/test/1.md", Question: "Q1", Answer: "A1", RevisitIn: 0}
	err := store.InsertFlashcard(fc)
	if err != nil {
		t.Fatalf("InsertFlashcard() error = %v", err)
	}

	// Get the inserted flashcard ID
	cards, err := store.GetAllFlashcards()
	if err != nil {
		t.Fatalf("GetAllFlashcards() error = %v", err)
	}

	if len(cards) != 1 {
		t.Fatalf("Expected 1 flashcard, got %d", len(cards))
	}

	id := cards[0].ID

	// Delete the flashcard
	err = store.DeleteFlashcard(id)
	if err != nil {
		t.Errorf("DeleteFlashcard() error = %v", err)
	}

	// Verify the deletion
	remainingCards, err := store.GetAllFlashcards()
	if err != nil {
		t.Fatalf("GetAllFlashcards() error = %v", err)
	}

	if len(remainingCards) != 0 {
		t.Errorf("Expected 0 flashcards after deletion, got %d", len(remainingCards))
	}
}

func TestIsFileProcessed(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	filePath := "/test/file.md"

	// Check before processing
	processed, err := store.IsFileProcessed(filePath)
	if err != nil {
		t.Fatalf("IsFileProcessed() error = %v", err)
	}

	if processed {
		t.Error("Expected file to be unprocessed initially")
	}

	// Insert a flashcard from the file
	fc := Flashcard{File: filePath, Question: "Q1", Answer: "A1", RevisitIn: 0}
	err = store.InsertFlashcard(fc)
	if err != nil {
		t.Fatalf("InsertFlashcard() error = %v", err)
	}

	// Check after processing
	processed, err = store.IsFileProcessed(filePath)
	if err != nil {
		t.Fatalf("IsFileProcessed() error = %v", err)
	}

	if !processed {
		t.Error("Expected file to be processed after insertion")
	}
}

func TestUpdateFlashcardFull(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	// Insert a flashcard
	fc := Flashcard{File: "/test/1.md", Question: "Q1", Answer: "A1", RevisitIn: 0}
	err := store.InsertFlashcard(fc)
	if err != nil {
		t.Fatalf("InsertFlashcard() error = %v", err)
	}

	// Get the inserted flashcard
	cards, err := store.GetAllFlashcards()
	if err != nil {
		t.Fatalf("GetAllFlashcards() error = %v", err)
	}

	if len(cards) != 1 {
		t.Fatalf("Expected 1 flashcard, got %d", len(cards))
	}

	// Update all fields
	cards[0].Question = "Q1 Updated"
	cards[0].Answer = "A1 Updated"
	cards[0].File = "/test/updated.md"
	cards[0].RevisitIn = 10
	err = store.UpdateFlashcardFull(cards[0])
	if err != nil {
		t.Errorf("UpdateFlashcardFull() error = %v", err)
	}

	// Verify the update
	updatedCards, err := store.GetAllFlashcards()
	if err != nil {
		t.Fatalf("GetAllFlashcards() error = %v", err)
	}

	if updatedCards[0].Question != "Q1 Updated" {
		t.Errorf("Expected Question 'Q1 Updated', got '%s'", updatedCards[0].Question)
	}
	if updatedCards[0].Answer != "A1 Updated" {
		t.Errorf("Expected Answer 'A1 Updated', got '%s'", updatedCards[0].Answer)
	}
	if updatedCards[0].File != "/test/updated.md" {
		t.Errorf("Expected File '/test/updated.md', got '%s'", updatedCards[0].File)
	}
	if updatedCards[0].RevisitIn != 10 {
		t.Errorf("Expected RevisitIn 10, got %d", updatedCards[0].RevisitIn)
	}
}

func TestGetUniqueFiles(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	// Insert flashcards from different files
	flashcards := []Flashcard{
		{File: "/test/file1.md", Question: "Q1", Answer: "A1", RevisitIn: 0},
		{File: "/test/file2.md", Question: "Q2", Answer: "A2", RevisitIn: 0},
		{File: "/test/file1.md", Question: "Q3", Answer: "A3", RevisitIn: 0}, // Duplicate file
		{File: "/test/file3.md", Question: "Q4", Answer: "A4", RevisitIn: 0},
	}

	for _, fc := range flashcards {
		err := store.InsertFlashcard(fc)
		if err != nil {
			t.Fatalf("InsertFlashcard() error = %v", err)
		}
	}

	// Get unique files
	files, err := store.GetUniqueFiles()
	if err != nil {
		t.Fatalf("GetUniqueFiles() error = %v", err)
	}

	// Should return 3 unique files
	if len(files) != 3 {
		t.Errorf("Expected 3 unique files, got %d", len(files))
	}

	// Check that files are sorted
	expectedFiles := []string{"/test/file1.md", "/test/file2.md", "/test/file3.md"}
	for i, expected := range expectedFiles {
		if i >= len(files) {
			t.Errorf("Missing file at index %d", i)
			continue
		}
		if files[i] != expected {
			t.Errorf("Expected file '%s' at index %d, got '%s'", expected, i, files[i])
		}
	}
}

func TestGetFlashcardsForReviewByFiles(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	// Insert flashcards from different files with different review statuses
	flashcards := []Flashcard{
		{File: "/test/file1.md", Question: "Q1", Answer: "A1", RevisitIn: 0},  // Due, file1
		{File: "/test/file1.md", Question: "Q2", Answer: "A2", RevisitIn: 5},  // Not due, file1
		{File: "/test/file2.md", Question: "Q3", Answer: "A3", RevisitIn: -1}, // Due, file2
		{File: "/test/file2.md", Question: "Q4", Answer: "A4", RevisitIn: 0},  // Due, file2
		{File: "/test/file3.md", Question: "Q5", Answer: "A5", RevisitIn: 0},  // Due, file3
	}

	for _, fc := range flashcards {
		err := store.InsertFlashcard(fc)
		if err != nil {
			t.Fatalf("InsertFlashcard() error = %v", err)
		}
	}

	// Test 1: Get flashcards for a single file
	selectedFiles := []string{"/test/file1.md"}
	cards, err := store.GetFlashcardsForReviewByFiles(selectedFiles)
	if err != nil {
		t.Fatalf("GetFlashcardsForReviewByFiles() error = %v", err)
	}

	if len(cards) != 1 {
		t.Errorf("Expected 1 flashcard for file1.md, got %d", len(cards))
	}

	if len(cards) > 0 && cards[0].Question != "Q1" {
		t.Errorf("Expected Q1, got %s", cards[0].Question)
	}

	// Test 2: Get flashcards for multiple files
	selectedFiles = []string{"/test/file1.md", "/test/file2.md"}
	cards, err = store.GetFlashcardsForReviewByFiles(selectedFiles)
	if err != nil {
		t.Fatalf("GetFlashcardsForReviewByFiles() error = %v", err)
	}

	if len(cards) != 3 {
		t.Errorf("Expected 3 flashcards for file1.md and file2.md, got %d", len(cards))
	}

	// Test 3: Get flashcards for all files
	selectedFiles = []string{"/test/file1.md", "/test/file2.md", "/test/file3.md"}
	cards, err = store.GetFlashcardsForReviewByFiles(selectedFiles)
	if err != nil {
		t.Fatalf("GetFlashcardsForReviewByFiles() error = %v", err)
	}

	if len(cards) != 4 {
		t.Errorf("Expected 4 flashcards for all files, got %d", len(cards))
	}

	// Test 4: Get flashcards for empty file list
	selectedFiles = []string{}
	cards, err = store.GetFlashcardsForReviewByFiles(selectedFiles)
	if err != nil {
		t.Fatalf("GetFlashcardsForReviewByFiles() error = %v", err)
	}

	if len(cards) != 0 {
		t.Errorf("Expected 0 flashcards for empty file list, got %d", len(cards))
	}

	// Test 5: Get flashcards for non-existent file
	selectedFiles = []string{"/test/nonexistent.md"}
	cards, err = store.GetFlashcardsForReviewByFiles(selectedFiles)
	if err != nil {
		t.Fatalf("GetFlashcardsForReviewByFiles() error = %v", err)
	}

	if len(cards) != 0 {
		t.Errorf("Expected 0 flashcards for non-existent file, got %d", len(cards))
	}
}

// setupTestDB creates a temporary database for testing
func setupTestDB(t *testing.T) *Store {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	store, err := NewStore(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	return store
}

func TestNewStore_ErrorHandling(t *testing.T) {
	// Test with invalid path (e.g., directory that doesn't exist and can't be created)
	invalidPath := "/invalid/path/that/does/not/exist/test.db"
	store, err := NewStore(invalidPath)

	// On some systems, this might succeed if the directory can be created
	// or fail if it can't. Either way, we should handle it gracefully.
	if err != nil {
		// Expected error for invalid path
		if store != nil {
			t.Error("NewStore() should return nil store on error")
		}
	} else {
		// If it succeeded, clean up
		if store != nil {
			store.Close()
		}
	}
}

func TestGetFlashcardsForReview_EmptyDatabase(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	cards, err := store.GetFlashcardsForReview()
	if err != nil {
		t.Fatalf("GetFlashcardsForReview() error = %v", err)
	}

	if len(cards) != 0 {
		t.Errorf("Expected 0 flashcards in empty database, got %d", len(cards))
	}
}

func TestGetUniqueFiles_EmptyDatabase(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	files, err := store.GetUniqueFiles()
	if err != nil {
		t.Fatalf("GetUniqueFiles() error = %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files in empty database, got %d", len(files))
	}
}

func TestGetAllFlashcards_EmptyDatabase(t *testing.T) {
	store := setupTestDB(t)
	defer store.Close()

	cards, err := store.GetAllFlashcards()
	if err != nil {
		t.Fatalf("GetAllFlashcards() error = %v", err)
	}

	if len(cards) != 0 {
		t.Errorf("Expected 0 flashcards in empty database, got %d", len(cards))
	}
}
