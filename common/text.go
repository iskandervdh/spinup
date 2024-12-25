package common

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Add a lipgloss info style to the given text.
func InfoText(text string) string {
	return fmt.Sprintln(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5DADE2")).
			Render(text),
	)
}

// Add a lipgloss success style to the given text.
func SuccessText(text string) string {
	return fmt.Sprintln(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A7E08F")).
			Render(text),
	)
}

// Add a lipgloss warning style to the given text.
func WarningText(text string) string {
	return fmt.Sprintln(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Render(text),
	)
}

// Add a lipgloss error style to the given text.
func ErrorText(text string) string {
	return fmt.Sprintln(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D90909")).
			Render("Error: " + text),
	)
}
