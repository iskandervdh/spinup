package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type selection struct {
	prompt  string
	options []string
	cursor  int
	exited  bool
}

func (s selection) Init() tea.Cmd {
	return nil
}

func (s selection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			s.exited = true
			return s, tea.Quit

		case "enter":
			return s, tea.Quit

		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}

		case "down", "j":
			if s.cursor < len(s.options)-1 {
				s.cursor++
			}
		}
	}

	return s, nil
}

func (s selection) View() string {
	out := fmt.Sprintf("%s\n\n", s.prompt)

	for i, choice := range s.options {
		cursor := " "
		if s.cursor == i {
			cursor = ">"
		}

		out += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	out += "\nPress enter to select.\n"

	return out
}
