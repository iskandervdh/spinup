package components

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/iskandervdh/spinup/common"
)

type Loading struct {
	spinner     spinner.Model
	loadingText string
	done        bool

	successText string
	errorText   string
}

func NewSpinner() spinner.Model {
	return spinner.New(
		spinner.WithSpinner(spinner.Dot),
		spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#a7e08f"))),
	)
}

func NewLoading(loadingText string) Loading {
	return Loading{
		spinner:     NewSpinner(),
		loadingText: loadingText,
	}
}

func (m Loading) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Loading) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	case common.SuccessMsg:
		m.done = true
		m.successText = msg.GetText()
		return m, tea.Quit

	case common.ErrMsg:
		m.done = true
		m.errorText = msg.GetText()
		return m, tea.Quit
	}

	return m, nil
}

func (m Loading) View() string {
	if m.errorText != "" {
		return common.ErrorText(m.errorText)
	}

	if m.done {
		return common.SuccessText(m.successText)
	}

	return fmt.Sprintf("\n    %s %s\n", m.spinner.View(), m.loadingText)
}
