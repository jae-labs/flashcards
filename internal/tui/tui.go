package tui

import (
	"flashcards/internal/tui/theme"
	"fmt"
)

// PrintInfo prints an informational message to the console
func PrintInfo(message string) {
	fmt.Println(theme.InfoStyle.Render(message))
}

// PrintError prints an error message to the console
func PrintError(message string, err error) {
	fmt.Println(theme.ErrorStyle.Render(message), err)
}

// PrintSuccess prints a success message to the console
func PrintSuccess(message string) {
	fmt.Println(theme.SuccessStyle.Render(message))
}
