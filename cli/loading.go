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
	doneText    string
	done        bool
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
	case spinner.TickMsg:
		if m.done {
			return m, tea.Quit
		}

		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case DoneMsg:
		m.done = true
		return m, tea.Quit
	}

	return m, nil
}

func (m loading) View() string {
	if m.done {
		return fmt.Sprintf(
			"%s\n",
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#a7e08f")).
				Render(m.doneText),
		)
	}

	return fmt.Sprintf("\n    %s %s\n", m.spinner.View(), m.loadingText)
}
