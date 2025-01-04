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

func (l Loading) Init() tea.Cmd {
	return l.spinner.Tick
}

func (l Loading) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			l.done = true
			return l, tea.Quit
		}

	case spinner.TickMsg:
		if l.done {
			return l, tea.Quit
		}

		var cmd tea.Cmd
		l.spinner, cmd = l.spinner.Update(msg)
		return l, cmd

	case *common.SuccessMsg:
		l.done = true
		l.successText = msg.GetText()
		return l, tea.Quit

	case *common.ErrMsg:
		l.done = true
		l.errorText = msg.GetText()
		return l, tea.Quit
	}

	return l, nil
}

func (l Loading) View() string {
	if l.errorText != "" {
		return common.ErrorText(l.errorText)
	}

	if l.done {
		return common.SuccessText(l.successText)
	}

	return fmt.Sprintf("\n    %s %s\n", l.spinner.View(), l.loadingText)
}

func (l Loading) GetErrorText() string {
	return l.errorText
}

func (l Loading) GetSuccessText() string {
	return l.successText
}
