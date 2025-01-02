package cli

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/iskandervdh/spinup/cli/components"
	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
	"github.com/iskandervdh/spinup/core"
)

// CLI struct that that determines the input and output of the CLI.
//
// It uses a message channel to communicate with the core.
type CLI struct {
	in  io.Reader
	out io.Writer

	core      *core.Core
	msgChan   *chan common.Msg
	msgChanWg *sync.WaitGroup
}

// Create a new CLI instance with the given options.
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

// Optional function to set the input of the CLI when creating a new instance.
func WithIn(in io.Reader) func(*CLI) {
	return func(c *CLI) {
		c.in = in
	}
}

// Optional function to set the output of the CLI when creating a new instance.
func WithOut(out io.Writer) func(*CLI) {
	return func(c *CLI) {
		c.out = out
	}
}

// Optional function to set the core of the CLI when creating a new instance.
func WithCore(core *core.Core) func(*CLI) {
	return func(c *CLI) {
		c.core = core
	}
}

// Clear the terminal screen by sending the escape codes to the output.
func (c *CLI) ClearTerminal() {
	fmt.Fprint(c.out, "\033[H\033[2J")
}

// Send a message to the message channel.
func (c *CLI) sendMsg(msg common.Msg) {
	*c.msgChan <- msg
}

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

	if _, err := p.Run(); err != nil {
		return common.NewErrMsg("Error starting program: %v", err)
	}

	return nil
}

func (c *CLI) sendHelpMsg() {
	c.sendMsg(common.NewRegularMsg("Usage: %s <command|project|variable|run|init> [args...]\n", config.ProgramName))
}

func (c *CLI) launchApp() {
	c.sendMsg(common.NewRegularMsg("Launching app...\n"))
	cmd := exec.Command(config.AppCommand)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if err != nil {
		panic(err)
	}

	err = cmd.Wait()

	if err != nil {
		panic(err)
	}
}

// Function to be called after the CLI has been initialized.
//
// It will handle the arguments passed to the CLI and
// execute the appropriate function based on the arguments.
func (c *CLI) Handle() {
	if len(os.Args) == 1 {
		if common.AppInstalled() {
			c.launchApp()
		} else {
			c.sendMsg(common.NewInfoMsg("App not installed. You can download it from https://github.com/iskandervdh/spinup-app/releases"))
			c.sendHelpMsg()
		}
	} else if os.Args[1] == "init" {
		c.sendMsg(c.core.Init())
	} else {
		c.core.FetchCommands()
		c.core.FetchProjects()

		switch os.Args[1] {
		case "-v", "--version":
			fmt.Printf("%s %s\n", config.ProgramName, strings.TrimSpace(config.Version))
		case "-h", "--help":
			c.sendHelpMsg()
		case "command", "c":
			c.handleCommand()
		case "project", "p":
			c.handleProject()
		case "variable", "v":
			c.handleVariable()
		case "domain-alias", "da":
			c.handleDomainAlias()
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
	}

	close(*c.msgChan)

	c.msgChanWg.Wait()
}
