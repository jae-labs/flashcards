package commands

import (
	"fmt"
	"os"

	"flashcards/internal/config"
	"flashcards/internal/store"
	"flashcards/internal/tui"

	"github.com/spf13/cobra"
)

var Store *store.Store
var Model string

var RootCmd = &cobra.Command{
	Use:   "flashcards",
	Short: "Ollama-powered spaced repetition flashcards CLI",
	Long: `Flashcards is a fast, minimal command-line tool for 
turning your notes into interactive flashcards and reviewing them with spaced 
repetition. Simply point Flashcards at your folder of markdown notes, and it uses 
Ollama's local AI models to automatically generate flashcards and quiz you in 
a colorful terminal interface.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg := config.LoadConfig()

		// Ensure data directory exists
		if err := cfg.EnsureDataDir(); err != nil {
			tui.PrintError("Could not create data directory:", err)
			os.Exit(1)
		}

		// Initialize database
		var err error
		Store, err = store.NewStore(cfg.DatabasePath)
		if err != nil {
			tui.PrintError("Database initialization failed:", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Default to review command
		ReviewCmd.Run(cmd, args)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if Store != nil {
			Store.Close()
		}
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVar(&Model, "model", "llama3.1", "Ollama model to use for flashcard generation")
	RootCmd.AddCommand(GenerateCmd)
	RootCmd.AddCommand(ReviewCmd)
	RootCmd.AddCommand(AdminCmd)
}
