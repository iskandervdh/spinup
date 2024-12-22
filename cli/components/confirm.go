package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Confirm struct {
	prompt string
	value  string
	exited bool
	cursor *cursor
}

func NewConfirm(prompt string) Confirm {
	return Confirm{
		prompt: prompt,
		value:  "",
		exited: false,
		cursor: newCursor(),
	}
}

func (c Confirm) GetValue() string {
	return c.value
}

func (c Confirm) GetExited() bool {
	return c.exited
}

func (c Confirm) Init() tea.Cmd {
	return c.cursor.init()
}

func (c Confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch strings.ToLower(msg.String()) {

		case "ctrl+c":
			c.exited = true
			return c, tea.Quit

		case "enter":
			return c, tea.Quit

		case "backspace":
			if c.cursor.position > 0 && len(c.value) > 0 {
				c.value = c.value[:c.cursor.position-1] + c.value[c.cursor.position:]
				c.cursor.moveLeft()
			}

		case "delete":
			if c.cursor.position < len(c.value) {
				c.value = c.value[:c.cursor.position] + c.value[c.cursor.position+1:]
			}

		case "left":
			c.cursor.moveLeft()

		case "right":
			c.cursor.moveRight(len(c.value))

		case "up", "down", "tab":
			break

		default:
			if c.cursor.position < len(c.value) {
				c.value = c.value[:c.cursor.position] + msg.String() + c.value[c.cursor.position:]
			} else {
				c.value += msg.String()
			}

			c.cursor.moveRight(len(c.value))
		}

	case blinkMsg:
		return c, c.cursor.toggle()
	}

	return c, nil
}

func (c Confirm) View() string {
	currentChar := " "

	if c.cursor.position < len(c.value) {
		currentChar = string(c.value[c.cursor.position])
	}

	valueWithCursor := c.value[:c.cursor.position] + c.cursor.get(currentChar)

	if c.cursor.position < len(c.value) {
		valueWithCursor += c.value[c.cursor.position+1:]
	}

	return c.prompt + " [Y/n] " + valueWithCursor
}
