package cli

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

func ErrorText(text string) string {
	return fmt.Sprintln(
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D90909")).
			Render("Error: " + text),
	)
}

func InfoPrint(a ...any) {
	fmt.Print(InfoText(fmt.Sprint(a...)))
}

func InfoPrintf(format string, a ...any) {
	fmt.Print(InfoText(fmt.Sprintf(format, a...)))
}

func SuccessPrint(a ...any) {
	fmt.Print(SuccessText(fmt.Sprint(a...)))
}

func SuccessPrintf(format string, a ...any) {
	fmt.Print(SuccessText(fmt.Sprintf(format, a...)))
}

func ErrorPrint(a ...any) {
	fmt.Print(ErrorText(fmt.Sprint(a...)))
}

func ErrorPrintf(format string, a ...any) {
	fmt.Print(ErrorText(fmt.Sprintf(format, a...)))
}
