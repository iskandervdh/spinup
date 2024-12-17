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

func (c *CLI) InfoPrint(a ...any) {
	fmt.Fprint(c.out, infoText(fmt.Sprint(a...)))
}

func (c *CLI) InfoPrintf(format string, a ...any) {
	fmt.Fprint(c.out, infoText(fmt.Sprintf(format, a...)))
}

func (c *CLI) SuccessPrint(a ...any) {
	fmt.Fprint(c.out, successText(fmt.Sprint(a...)))
}

func (c *CLI) SuccessPrintf(format string, a ...any) {
	fmt.Fprint(c.out, successText(fmt.Sprintf(format, a...)))
}

func (c *CLI) WarningPrint(a ...any) {
	fmt.Fprint(c.out, warningText(fmt.Sprint(a...)))
}

func (c *CLI) WarningPrintf(format string, a ...any) {
	fmt.Fprint(c.out, warningText(fmt.Sprintf(format, a...)))
}

func (c *CLI) ErrorPrint(a ...any) {
	fmt.Fprint(c.out, errorText(fmt.Sprint(a...)))
}

func (c *CLI) ErrorPrintf(format string, a ...any) {
	fmt.Fprint(c.out, errorText(fmt.Sprintf(format, a...)))
}
