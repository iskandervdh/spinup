package components

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type blinkMsg struct{}

func blink() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(time.Time) tea.Msg {
		return blinkMsg{}
	})
}

type cursor struct {
	visible  bool
	position int
}

func newCursor() *cursor {
	return &cursor{visible: true, position: 0}
}

func (c *cursor) init() tea.Cmd {
	c.visible = true

	return blink()
}

func (c *cursor) toggle() tea.Cmd {
	c.visible = !c.visible

	return blink()
}

func (c *cursor) moveLeft() {
	if c.position > 0 {
		c.position--
		c.visible = true
	}
}

func (c *cursor) moveRight(max int) {
	if c.position < max {
		c.position++
		c.visible = true
	}
}

func (c *cursor) get(currentChar string) string {
	if c.visible {
		return "\x1b[7m" + currentChar + "\x1b[0m"
	}

	return currentChar
}
