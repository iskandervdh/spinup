package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/iskandervdh/spinup/cli/components"
	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
	"github.com/iskandervdh/spinup/core"
)

type CLI struct {
	in  io.Reader
	out io.Writer

	core      *core.Core
	msgChan   *chan common.Msg
	msgChanWg *sync.WaitGroup
}

func New(options ...func(*CLI)) *CLI {
	msgChan := make(chan common.Msg, 100)
	msgChanWg := sync.WaitGroup{}

	c := &CLI{
		in:        os.Stdin,
		out:       os.Stdout,
		core:      core.New(core.WithMsgChan(&msgChan)),
		msgChan:   &msgChan,
		msgChanWg: &msgChanWg,
	}

	for _, option := range options {
		option(c)
	}

	c.msgChanWg.Add(1)

	go func() {
		defer msgChanWg.Done()

		for msg := range *c.msgChan {
			c.MsgPrint(msg)
		}
	}()

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

func WithCore(core *core.Core) func(*CLI) {
	return func(c *CLI) {
		c.core = core
	}
}

func (c *CLI) ClearTerminal() {
	fmt.Fprint(c.out, "\033[H\033[2J")
}

func (c *CLI) sendMsg(msg common.Msg) {
	*c.msgChan <- msg
}

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

func (c *CLI) Loading(loadingText string, f func() common.Msg) common.Msg {
	l := components.NewLoading(loadingText)

	p := tea.NewProgram(l)

	go func() {
		msg := f()
		p.Send(msg)
	}()

	if _, err := p.Run(); err != nil {
		return common.NewErrMsg("Error starting program: %v", err)
	}

	return nil
}

func (c *CLI) Handle() {
	if len(os.Args) < 2 {
		c.sendMsg(common.NewRegularMsg("Usage: %s <command|project|variable|run|init> [args...]\n", config.ProgramName))
		return
	}

	if os.Args[1] == "init" {
		c.sendMsg(c.core.Init())
		return
	}

	c.core.GetCommandsConfig()
	c.core.GetProjectsConfig()

	switch os.Args[1] {
	case "-v", "--version":
		fmt.Printf("%s %s\n", config.ProgramName, strings.TrimSpace(config.Version))
	case "command", "c":
		c.handleCommand()
	case "project", "p":
		c.handleProject()
	case "variable", "v":
		c.handleVariable()
	case "run":
		if len(os.Args) < 3 {
			c.sendMsg(common.NewRegularMsg("Usage: %s run <project>\n", config.ProgramName))
			break
		}

		result := c.core.TryToRun(os.Args[2])

		if _, ok := result.(*common.ErrMsg); ok {
			c.ErrorPrint(result)
			os.Exit(1)
		}

		if result == nil {
			c.sendMsg(common.NewErrMsg("Unknown project '%s'\n", os.Args[2]))
		}
	default:
		result := c.core.TryToRun(os.Args[1])

		if _, ok := result.(*common.ErrMsg); ok {
			c.ErrorPrint(result)
			os.Exit(1)
		}

		if result == nil {
			c.sendMsg(common.NewRegularMsg("Unknown subcommand or project '%s'\n\n", os.Args[1]))
			c.sendMsg(common.NewRegularMsg("Expected 'command|c', 'project|p', 'run' or 'init' subcommand or a valid project name\n"))
		}
	}

	close(*c.msgChan)

	c.msgChanWg.Wait()
}
