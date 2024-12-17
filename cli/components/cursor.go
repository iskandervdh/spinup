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
	visible bool
}

func newCursor() *cursor {
	return &cursor{visible: true}
}

func (c *cursor) init() tea.Cmd {
	c.visible = true
	return blink()
}

func (c *cursor) toggle() tea.Cmd {
	c.visible = !c.visible

	return blink()
}

func (c *cursor) get() string {
	if c.visible {
		return "\x1b[7m \x1b[0m"
	}

	return ""
}
