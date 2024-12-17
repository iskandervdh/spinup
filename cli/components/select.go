package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Selection struct {
	prompt  string
	options []string
	cursor  int
	exited  bool
}

func NewSelection(prompt string, options []string) Selection {
	return Selection{
		prompt:  prompt,
		options: options,
		cursor:  0,
		exited:  false,
	}
}

func (s Selection) GetValue() string {
	return s.options[s.cursor]
}

func (s Selection) GetExited() bool {
	return s.exited
}

func (s Selection) Init() tea.Cmd {
	return nil
}

func (s Selection) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (s Selection) View() string {
	out := fmt.Sprintf("\n%s\n\n", s.prompt)

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
