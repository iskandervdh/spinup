package common

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func InfoText(text string) string {
	return fmt.Sprintln(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5DADE2")).
			Render(text),
	)
}

func SuccessText(text string) string {
	return fmt.Sprintln(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A7E08F")).
			Render(text),
	)
}

func WarningText(text string) string {
	return fmt.Sprintln(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Render(text),
	)
}

func ErrorText(text string) string {
	return fmt.Sprintln(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D90909")).
			Render("Error: " + text),
	)
}
