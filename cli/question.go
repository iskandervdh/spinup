package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type question struct {
	prompt   string
	options  []string
	cursor   int
	selected []bool
	exited   bool
}

func (q question) Init() tea.Cmd {
	return nil
}

func (q question) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			q.exited = true
			return q, tea.Quit

		case "enter":
			return q, tea.Quit

		case "up", "k":
			if q.cursor > 0 {
				q.cursor--
			}

		case "down", "j":
			if q.cursor < len(q.options)-1 {
				q.cursor++
			}

		case " ":
			checked := q.selected[q.cursor]
			if checked {
				q.selected[q.cursor] = false
			} else {
				q.selected[q.cursor] = true
			}
		}
	}

	return q, nil
}

func (q question) View() string {
	s := fmt.Sprintf("%s\n\n", q.prompt)

	for i, choice := range q.options {
		cursor := " "
		if q.cursor == i {
			cursor = ">"
		}

		checked := " "

		if c := q.selected[i]; c {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress enter to submit.\n"

	return s
}