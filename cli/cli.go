package cli

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func ClearTerminal() {
	fmt.Print("\033[H\033[2J")
}

func Question(prompt string, options []string) []string {
	q := question{
		prompt:   prompt,
		options:  options,
		selected: make([]bool, len(options)),
		exited:   false,
	}

	p := tea.NewProgram(q)
	m, err := p.Run()

	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	r := m.(question)

	if r.exited {
		os.Exit(0)
	}

	selectedOptions := make([]string, 0, len(q.selected))

	for i, checked := range r.selected {
		if checked {
			selectedOptions = append(selectedOptions, q.options[i])
		}
	}

	return selectedOptions
}

func Selection(prompt string, options []string) string {
	s := selection{
		prompt:  prompt,
		options: options,
		exited:  false,
	}

	p := tea.NewProgram(s)
	m, err := p.Run()

	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	r := m.(selection)

	if r.exited {
		os.Exit(0)
	}

	return r.options[r.cursor]
}

func Input(prompt string) string {
	i := input{
		prompt: prompt,
		value:  "",
		exited: false,
		cursor: newCursor(),
	}

	p := tea.NewProgram(i)
	m, err := p.Run()

	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	r := m.(input)

	if r.exited {
		os.Exit(0)
	}

	return r.value
}

func Confirm(prompt string) bool {
	c := confirm{
		prompt: prompt,
		value:  "",
		exited: false,
		cursor: newCursor(),
	}

	p := tea.NewProgram(c)
	m, err := p.Run()

	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	r := m.(confirm)

	if r.exited {
		os.Exit(0)
	}

	switch strings.ToLower(r.value) {
	case "y", "yes":
		return true
	}

	return false
}
