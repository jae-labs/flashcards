package commands

import (
	"fmt"

	"flashcards/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Flashcards database management",
	Run: func(cmd *cobra.Command, args []string) {
		// Store is initialized in RootCmd PersistentPreRun
		list, err := Store.GetAllFlashcards()
		if err != nil {
			fmt.Println("DB error:", err)
			return
		}
		model := tui.NewAdminModel(Store, list)
		if _, err := tea.NewProgram(model).Run(); err != nil {
			fmt.Println("Error running admin TUI:", err)
		}
	},
}
