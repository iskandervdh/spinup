package components

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Input struct {
	prompt string
	value  string
	exited bool
	cursor *cursor
}

func NewInput(prompt string) Input {
	return Input{
		prompt: prompt,
		value:  "",
		exited: false,
		cursor: newCursor(),
	}
}

func (i Input) GetValue() string {
	return i.value
}

func (i Input) GetExited() bool {
	return i.exited
}

func (i Input) Init() tea.Cmd {
	return i.cursor.init()
}

func (i Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (i Input) View() string {
	return i.prompt + " " + i.value + i.cursor.get()
}
