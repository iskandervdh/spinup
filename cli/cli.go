package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type doneMsg struct {
	text string
}

type errMsg struct {
	text string
}

type CLI struct {
	in  io.Reader
	out io.Writer
}

func DoneMsg(text string) doneMsg {
	return doneMsg{text: text}
}

func ErrMsg(text string) errMsg {
	return errMsg{text: text}
}

func New(options ...func(*CLI)) *CLI {
	c := &CLI{
		in:  os.Stdin,
		out: os.Stdout,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func WithIn(in io.Reader) func(*CLI) {
	return func(c *CLI) {
		c.in = in
	}
}

func WithOut(out io.Writer) func(*CLI) {
	return func(c *CLI) {
		c.out = out
	}
}

func (c *CLI) ClearTerminal() {
	fmt.Fprint(c.out, "\033[H\033[2J")
}

func (c *CLI) Question(prompt string, options []string) ([]string, error, bool) {
	q := question{
		prompt:   prompt,
		options:  options,
		selected: make([]bool, len(options)),
		exited:   false,
	}

	p := tea.NewProgram(q, tea.WithInput(c.in), tea.WithOutput(c.out))

	m, err := p.Run()

	if err != nil {
		return nil, err, false
	}

	r := m.(question)

	if r.exited {
		return nil, nil, r.exited
	}

	selectedOptions := make([]string, 0, len(q.selected))

	for i, checked := range r.selected {
		if checked {
			selectedOptions = append(selectedOptions, q.options[i])
		}
	}

	return selectedOptions, nil, r.exited
}

func (c *CLI) Selection(prompt string, options []string) (string, error, bool) {
	s := selection{
		prompt:  prompt,
		options: options,
		exited:  false,
	}

	p := tea.NewProgram(s, tea.WithInput(c.in), tea.WithOutput(c.out))
	m, err := p.Run()

	if err != nil {
		return "", err, false
	}

	r := m.(selection)

	if r.exited {
		return "", nil, r.exited
	}

	return r.options[r.cursor], nil, r.exited
}

func (c *CLI) Input(prompt string) string {
	i := input{
		prompt: prompt,
		value:  "",
		exited: false,
		cursor: newCursor(),
	}

	p := tea.NewProgram(i, tea.WithInput(c.in), tea.WithOutput(c.out))
	m, err := p.Run()

	if err != nil {
		c.ErrorPrint(err)
		os.Exit(1)
	}

	r := m.(input)

	if r.exited {
		os.Exit(0)
	}

	return r.value
}

func (c *CLI) Confirm(prompt string) bool {
	conf := confirm{
		prompt: prompt,
		value:  "",
		exited: false,
		cursor: newCursor(),
	}

	p := tea.NewProgram(conf, tea.WithInput(c.in), tea.WithOutput(c.out))
	m, err := p.Run()

	if err != nil {
		c.ErrorPrint(err)
		os.Exit(1)
	}

	r := m.(confirm)

	if r.exited {
		os.Exit(0)
	}

	switch strings.ToLower(r.value) {
	case "y", "yes":
		return true
	}

	return false
}

func (c *CLI) Loading(loadingText string, f func() tea.Msg) *loading {
	s := newLoadingSpinner()

	l := loading{
		spinner:     s,
		loadingText: loadingText,
	}

	p := tea.NewProgram(l)

	go func() {
		msg := f()
		p.Send(msg)
	}()

	m, err := p.Run()

	if err != nil {
		c.ErrorPrintf("Error starting program: %v", err)
		return nil
	}

	r := m.(loading)

	return &r
}
