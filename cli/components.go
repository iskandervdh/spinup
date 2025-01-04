package cli

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/iskandervdh/spinup/cli/components"
	"github.com/iskandervdh/spinup/common"
)

// CLI handling of Question component.
func (c *CLI) Question(prompt string, options []string, defaultSelected []bool) ([]string, error, bool) {
	q := components.NewQuestion(prompt, options, defaultSelected)

	p := tea.NewProgram(q, tea.WithInput(c.in), tea.WithOutput(c.out))

	m, err := p.Run()

	if err != nil {
		return nil, err, false
	}

	r := m.(components.Question)

	if r.GetExited() {
		return nil, nil, true
	}

	return r.GetSelected(), nil, false
}

// CLI handling of Selection component.
func (c *CLI) Selection(prompt string, options []string) (string, error, bool) {
	s := components.NewSelection(prompt, options)

	p := tea.NewProgram(s, tea.WithInput(c.in), tea.WithOutput(c.out))
	m, err := p.Run()

	if err != nil {
		return "", err, false
	}

	r := m.(components.Selection)

	if r.GetExited() {
		return "", nil, true
	}

	return r.GetValue(), nil, false
}

// CLI handling of Input component.
func (c *CLI) Input(prompt string, defaultValue string) string {
	i := components.NewInput(prompt, defaultValue)

	p := tea.NewProgram(i, tea.WithInput(c.in), tea.WithOutput(c.out))
	m, err := p.Run()

	if err != nil {
		c.ErrorPrint(err)
		os.Exit(1)
	}

	r := m.(components.Input)

	if r.GetExited() {
		os.Exit(0)
	}

	return r.GetValue()
}

// CLI handling of Confirm component.
func (c *CLI) Confirm(prompt string) bool {
	conf := components.NewConfirm(prompt)

	p := tea.NewProgram(conf, tea.WithInput(c.in), tea.WithOutput(c.out))
	m, err := p.Run()

	if err != nil {
		c.ErrorPrint(err)
		os.Exit(1)
	}

	r := m.(components.Confirm)

	if r.GetExited() {
		os.Exit(0)
	}

	switch strings.ToLower(r.GetValue()) {
	case "y", "yes":
		return true
	}

	return false
}

// CLI handling of Loading component.
func (c *CLI) Loading(loadingText string, f func() common.Msg) common.Msg {
	l := components.NewLoading(loadingText)

	p := tea.NewProgram(l)

	go func() {
		msg := f()
		p.Send(msg)
	}()

	m, err := p.Run()

	if err != nil {
		return common.NewErrMsg("Error starting program: %v", err)
	}

	loading := m.(components.Loading)

	if loading.GetSuccessText() != "" {
		return common.NewSuccessMsg(l.GetSuccessText())
	}

	return common.NewErrMsg(l.GetErrorText())

}
