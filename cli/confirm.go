package cli

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type confirm struct {
	prompt string
	value  string
	exited bool
	cursor *cursor
}

func (c confirm) Init() tea.Cmd {
	return c.cursor.init()
}

func (c confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {

		case "ctrl+c":
			c.exited = true
			return c, tea.Quit

		case "enter":
			return c, tea.Quit

		case "backspace":
			if len(c.value) > 0 {
				c.value = c.value[:len(c.value)-1]
			}

		default:
			c.value += msg.String()
		}

	case blinkMsg:
		return c, c.cursor.toggle()
	}

	return c, nil
}

func (c confirm) View() string {
	return c.prompt + " [Y/n] " + c.value + c.cursor.get()
}
