package commands

import (
	"flashcards/internal/tui"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var ReviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review flashcards",
	Run: func(cmd *cobra.Command, args []string) {
		// Step 1: Get all unique files from database
		allFiles, err := Store.GetUniqueFiles()
		if err != nil {
			tui.PrintError("Failed to get files from database:", err)
			return
		}

		if len(allFiles) == 0 {
			tui.PrintInfo("No files found in the database. Generate flashcards first.")
			return
		}

		// Step 2: Show file selector UI
		fileSelector := tui.NewFileSelectorModel(allFiles)
		p := tea.NewProgram(fileSelector)
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running file selector:", err)
			return
		}

		// Step 3: Get selected files
		selectedFiles := fileSelector.GetSelectedFiles()
		if len(selectedFiles) == 0 {
			tui.PrintInfo("No files selected. See you next time!")
			return
		}

		// Step 4: Get flashcards for selected files
		flashcards, err := Store.GetFlashcardsForReviewByFiles(selectedFiles)
		if err != nil {
			tui.PrintError("DB query error:", err)
			return
		}

		if len(flashcards) == 0 {
			tui.PrintInfo("No flashcards due for review in the selected file(s). Well done!")
			return
		}

		// Step 5: Run Bubble Tea TUI for review
		model := tui.NewReviewModel(flashcards)
		p = tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running review TUI:", err)
		}

		// Step 6: After review, update DB only for flashcards that were actually answered
		for i, fc := range flashcards {
			if model.FlashcardWasCorrect(i) && model.FlashcardRevisitIn(i) > 0 {
				fc.RevisitIn = model.FlashcardRevisitIn(i)
				if err := Store.UpdateFlashcard(fc); err != nil {
					tui.PrintError("DB update error:", err)
				} else {
					tui.PrintSuccess(fmt.Sprintf("Updated flashcard %d: revisitin=%d", fc.ID, fc.RevisitIn))
				}
			} else if !model.FlashcardWasCorrect(i) && model.FlashcardRevisitIn(i) > 0 {
				// Incorrect -> revisit sooner (e.g. tomorrow => 1)
				fc.RevisitIn = 1
				if err := Store.UpdateFlashcard(fc); err != nil {
					tui.PrintError("DB update error:", err)
				} else {
					tui.PrintSuccess(fmt.Sprintf("Marked flashcard %d incorrect: revisitin set to 1", fc.ID))
				}
			}
		}
	},
}
