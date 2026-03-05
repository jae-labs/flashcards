// Package store provides data persistence for flashcards using SQLite
package store

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Store manages the database connection and operations for flashcards
type Store struct {
	DB *sql.DB // SQLite database connection
}

// NewStore creates a new Store instance with the specified database file
// It automatically creates the flashcards table if it doesn't exist
func NewStore(dbName string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	// Configure connection pool for better performance
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Create flashcards table with proper schema
	createTable := `CREATE TABLE IF NOT EXISTS flashcards (
			  id INTEGER PRIMARY KEY AUTOINCREMENT,
			  file TEXT NOT NULL,
			  question TEXT NOT NULL,
			  answer TEXT NOT NULL,
			  revisitin INTEGER DEFAULT 0,
			  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		  );`
	_, err = db.Exec(createTable)
	if err != nil {
		return nil, err
	}

	// Create indexes for frequently queried columns to improve performance
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_flashcards_revisitin ON flashcards(revisitin)`,
		`CREATE INDEX IF NOT EXISTS idx_flashcards_file ON flashcards(file)`,
		`CREATE INDEX IF NOT EXISTS idx_flashcards_file_revisitin ON flashcards(file, revisitin)`,
	}
	for _, idx := range indexes {
		if _, err := db.Exec(idx); err != nil {
			return nil, fmt.Errorf("failed to create index: %w", err)
		}
	}

	return &Store{DB: db}, nil
}

// GetFlashcardsForReview returns all flashcards that are due for review
// A flashcard is due for review when RevisitIn <= 0 or when the revisit date has passed
func (s *Store) GetFlashcardsForReview() ([]Flashcard, error) {
	query := `SELECT id, file, question, answer, revisitin 
			  FROM flashcards 
			  WHERE revisitin <= 0 
			  ORDER BY id ASC`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query flashcards for review: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	// Pre-allocate slice with reasonable initial capacity to reduce allocations
	flashcards := make([]Flashcard, 0, 100)
	for rows.Next() {
		var fc Flashcard
		err := rows.Scan(&fc.ID, &fc.File, &fc.Question, &fc.Answer, &fc.RevisitIn)
		if err != nil {
			return nil, fmt.Errorf("failed to scan flashcard: %w", err)
		}
		flashcards = append(flashcards, fc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating flashcards: %w", err)
	}

	return flashcards, nil
}

// GetFlashcardsForReviewByFiles returns flashcards due for review filtered by specific file paths
func (s *Store) GetFlashcardsForReviewByFiles(files []string) ([]Flashcard, error) {
	if len(files) == 0 {
		return []Flashcard{}, nil
	}

	// Build placeholders for IN clause
	placeholders := ""
	for i := range files {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "?"
	}

	// nosemgrep: go.lang.security.audit.database.string-formatted-query.string-formatted-query
	// #nosec G201 -- This is safe: we're only using fmt.Sprintf to build placeholders (?), not user data
	query := fmt.Sprintf(`SELECT id, file, question, answer, revisitin 
			  FROM flashcards 
			  WHERE revisitin <= 0 AND file IN (%s)
			  ORDER BY id ASC`, placeholders)

	// Convert files to []interface{} for Query
	args := make([]interface{}, len(files))
	for i, f := range files {
		args[i] = f
	}

	rows, err := s.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query flashcards for review by files: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	// Pre-allocate slice with reasonable initial capacity to reduce allocations
	flashcards := make([]Flashcard, 0, 50)
	for rows.Next() {
		var fc Flashcard
		err := rows.Scan(&fc.ID, &fc.File, &fc.Question, &fc.Answer, &fc.RevisitIn)
		if err != nil {
			return nil, fmt.Errorf("failed to scan flashcard: %w", err)
		}
		flashcards = append(flashcards, fc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating flashcards: %w", err)
	}

	return flashcards, nil
}

// GetUniqueFiles returns all unique file paths that have flashcards in the database
func (s *Store) GetUniqueFiles() ([]string, error) {
	query := `SELECT DISTINCT file FROM flashcards ORDER BY file ASC`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query unique files: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	// Pre-allocate slice with reasonable initial capacity to reduce allocations
	files := make([]string, 0, 20)
	for rows.Next() {
		var file string
		err := rows.Scan(&file)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}
		files = append(files, file)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating files: %w", err)
	}

	return files, nil
}

// GetAllFlashcards returns all flashcards ordered by revisitin ascending
func (s *Store) GetAllFlashcards() ([]Flashcard, error) {
	rows, err := s.DB.Query("SELECT id, file, question, answer, revisitin FROM flashcards ORDER BY revisitin ASC, id ASC")
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	// Pre-allocate slice with reasonable initial capacity to reduce allocations
	flashcards := make([]Flashcard, 0, 100)
	for rows.Next() {
		var fc Flashcard
		if err := rows.Scan(&fc.ID, &fc.File, &fc.Question, &fc.Answer, &fc.RevisitIn); err != nil {
			return nil, err
		}
		flashcards = append(flashcards, fc)
	}
	return flashcards, nil
}

// DeleteFlashcard deletes a flashcard by id
func (s *Store) DeleteFlashcard(id int) error {
	_, err := s.DB.Exec("DELETE FROM flashcards WHERE id=?", id)
	return err
}

// UpdateFlashcardFull updates all editable fields of a flashcard
func (s *Store) UpdateFlashcardFull(fc Flashcard) error {
	_, err := s.DB.Exec("UPDATE flashcards SET file=?, question=?, answer=?, revisitin=? WHERE id=?", fc.File, fc.Question, fc.Answer, fc.RevisitIn, fc.ID)
	return err
}

// UpdateFlashcard updates a flashcard's revisitin date
func (s *Store) UpdateFlashcard(fc Flashcard) error {
	_, err := s.DB.Exec("UPDATE flashcards SET revisitin=? WHERE id=?", fc.RevisitIn, fc.ID)
	return err
}

// IsFileProcessed checks if a file has already been processed
func (s *Store) IsFileProcessed(filePath string) (bool, error) {
	var count int
	row := s.DB.QueryRow("SELECT COUNT(*) FROM flashcards WHERE file = ?", filePath)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// InsertFlashcard inserts a new flashcard into the database
func (s *Store) InsertFlashcard(fc Flashcard) error {
	_, err := s.DB.Exec("INSERT INTO flashcards (file, question, answer, revisitin) VALUES (?, ?, ?, ?)", fc.File, fc.Question, fc.Answer, fc.RevisitIn)
	return err
}

// Close closes the database connection
func (s *Store) Close() {
	_ = s.DB.Close()
}
