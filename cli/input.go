package cli

import (
	tea "github.com/charmbracelet/bubbletea"
)

type input struct {
	prompt string
	value  string
	exited bool
	cursor *cursor
}

func (i input) Init() tea.Cmd {
	return i.cursor.init()
}

func (i input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c":
			i.exited = true
			return i, tea.Quit

		case "enter":
			return i, tea.Quit

		case "backspace":
			if len(i.value) > 0 {
				i.value = i.value[:len(i.value)-1]
			}

		default:
			i.value += msg.String()
		}

	case blinkMsg:
		return i, i.cursor.toggle()
	}

	return i, nil
}

func (i input) View() string {
	return i.prompt + " " + i.value + i.cursor.get()
}
