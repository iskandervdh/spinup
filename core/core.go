package core

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
)

// Core is the main struct for the core package.
//
// It contains the config, message channel, commands and projects.
// The config is used to read and update the configuration.
// The message channel is used to send messages to the CLI or the app.
type Core struct {
	config  *config.Config
	msgChan *chan common.Msg
	sigChan *chan os.Signal

	out io.Writer
	err io.Writer

	commands Commands
	projects Projects
}

// Create a new Core instance with the given options.
func New(options ...func(*Core)) *Core {
	config, err := config.New()

	if err != nil {
		fmt.Println("Error getting config:", err)
		os.Exit(1)
	}

	s := &Core{
		config: config,
		out:    os.Stdout,
		err:    os.Stderr,
	}

	for _, option := range options {
		option(s)
	}

	return s
}

// Optional function to set the config of the Core when creating a new instance.
func WithConfig(config *config.Config) func(*Core) {
	return func(c *Core) {
		c.config = config
	}
}

// Optional function to set the message channel of the Core when creating a new instance.
func WithMsgChan(msgChan *chan common.Msg) func(*Core) {
	return func(c *Core) {
		c.msgChan = msgChan
	}
}

func (c *Core) SetOut(out io.Writer) {
	c.out = out
}

func (c *Core) SetErr(err io.Writer) {
	c.err = err
}

// Get the config of the Core instance.
func (c *Core) GetConfig() *config.Config {
	return c.config
}

// RequireSudo checks if the user has sudo permissions.
// This is needed to update some system files.
//
// It returns an error if the user does not have sudo permissions.
func (c *Core) RequireSudo() error {
	if c.config.IsTesting() {
		return nil
	}

	err := exec.Command("sudo", "-v").Run()

	if err != nil {
		return fmt.Errorf("this command requires sudo")
	}

	return nil
}

// Get the commands from the commands.json file.
func (c *Core) GetCommandsConfig() error {
	commands, err := c.GetCommands()

	if err != nil {
		return fmt.Errorf("error getting commands. Did you run init?")
	}

	c.commands = commands

	return nil
}

// Get the projects from the projects.json file.
func (c *Core) GetProjectsConfig() error {
	projects, err := c.GetProjects()

	if err != nil {
		return fmt.Errorf("error getting projects. Did you run init?")
	}

	c.projects = projects

	return nil
}

// Get all the names of the commands.
func (c *Core) GetCommandNames() []string {
	if c.commands == nil {
		c.GetCommandsConfig()
	}

	var commandNames []string

	for commandName := range c.commands {
		commandNames = append(commandNames, commandName)
	}

	return commandNames
}

// Get all the names of the projects.
func (c *Core) GetProjectNames() []string {
	if c.projects == nil {
		c.GetProjectsConfig()
	}

	var projectNames []string

	for commandName := range c.projects {
		projectNames = append(projectNames, commandName)
	}

	return projectNames
}

// Get the signal channel of the Core instance.
func (c *Core) GetSigChan() *chan os.Signal {
	return c.sigChan
}

// Get all the commands that are part of the given project.
func (c *Core) getCommandsForProject(projectName string) []string {
	if c.projects == nil {
		c.GetProjectsConfig()
	}

	project, ok := c.projects[projectName]

	if !ok {
		return nil
	}

	return project.Commands
}

// Send a message to the message channel.
func (c *Core) sendMsg(msg common.Msg) {
	*c.msgChan <- msg
}
