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

func NewInput(prompt string, defaultValue string) Input {
	input := Input{
		prompt: prompt,
		value:  defaultValue,
		exited: false,
		cursor: newCursor(),
	}

	if defaultValue != "" {
		input.cursor.position = len(defaultValue)
	}

	return input
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
			if i.cursor.position > 0 && len(i.value) > 0 {
				i.value = i.value[:i.cursor.position-1] + i.value[i.cursor.position:]
				i.cursor.moveLeft()
			}

		case "delete":
			if i.cursor.position < len(i.value) {
				i.value = i.value[:i.cursor.position] + i.value[i.cursor.position+1:]
			}

		case "left":
			i.cursor.moveLeft()

		case "right":
			i.cursor.moveRight(len(i.value))

		case "up", "down", "tab":
			break

		default:
			if i.cursor.position < len(i.value) {
				i.value = i.value[:i.cursor.position] + msg.String() + i.value[i.cursor.position:]
			} else {
				i.value += msg.String()
			}

			i.cursor.moveRight(len(i.value))
		}

	case blinkMsg:
		return i, i.cursor.toggle()
	}

	return i, nil
}

func (i Input) View() string {
	currentChar := " "

	if i.cursor.position < len(i.value) {
		currentChar = string(i.value[i.cursor.position])
	}

	valueWithCursor := i.value[:i.cursor.position] + i.cursor.get(currentChar)

	if i.cursor.position < len(i.value) {
		valueWithCursor += i.value[i.cursor.position+1:]
	}

	return i.prompt + " " + valueWithCursor
}
