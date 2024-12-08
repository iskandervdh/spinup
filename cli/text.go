package cli

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func infoText(text string) string {
	return fmt.Sprintln(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5DADE2")).
			Render(text),
	)
}

func successText(text string) string {
	return fmt.Sprintln(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A7E08F")).
			Render(text),
	)
}

func warningText(text string) string {
	return fmt.Sprintln(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Render(text),
	)
}

func errorText(text string) string {
	return fmt.Sprintln(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D90909")).
			Render("Error: " + text),
	)
}

func InfoPrint(a ...any) {
	fmt.Print(infoText(fmt.Sprint(a...)))
}

func InfoPrintf(format string, a ...any) {
	fmt.Print(infoText(fmt.Sprintf(format, a...)))
}

func SuccessPrint(a ...any) {
	fmt.Print(successText(fmt.Sprint(a...)))
}

func SuccessPrintf(format string, a ...any) {
	fmt.Print(successText(fmt.Sprintf(format, a...)))
}

func WarningPrint(a ...any) {
	fmt.Print(warningText(fmt.Sprint(a...)))
}

func WarningPrintf(format string, a ...any) {
	fmt.Print(warningText(fmt.Sprintf(format, a...)))
}

func ErrorPrint(a ...any) {
	fmt.Print(errorText(fmt.Sprint(a...)))
}

func ErrorPrintf(format string, a ...any) {
	fmt.Print(errorText(fmt.Sprintf(format, a...)))
}
