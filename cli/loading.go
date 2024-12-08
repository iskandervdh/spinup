package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type loading struct {
	spinner     spinner.Model
	loadingText string
	done        bool

	doneText  string
	errorText string
}

func newLoadingSpinner() spinner.Model {
	return spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#a7e08f"))),
	)
}

func (m loading) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m loading) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case spinner.TickMsg:
		if m.done {
			return m, tea.Quit
		}

		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case doneMsg:
		m.done = true
		m.doneText = msg.text
		return m, tea.Quit

	case errMsg:
		m.done = true
		m.errorText = msg.text
		return m, tea.Quit
	}

	return m, nil
}

func (m loading) View() string {
	if m.errorText != "" {
		return errorText(m.errorText)
	}

	if m.done {
		return successText(m.doneText)
	}

	return fmt.Sprintf("\n    %s %s\n", m.spinner.View(), m.loadingText)
}
