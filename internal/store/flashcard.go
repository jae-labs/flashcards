// Package store provides data persistence for flashcards using SQLite
package store

// Flashcard represents a single flashcard with spaced repetition metadata
type Flashcard struct {
	ID        int    // Unique identifier for the flashcard
	File      string // Source file path where the flashcard was generated from
	Question  string // The question/text to be reviewed
	Answer    string // The answer/explanation for the question
	RevisitIn int    // Number of days until next review (<=0 means due for review)
}
