package core

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
)

type Core struct {
	config  *config.Config
	msgChan *chan common.Msg

	commands Commands
	projects Projects
}

func New(options ...func(*Core)) *Core {
	config, err := config.New()

	if err != nil {
		fmt.Println("Error getting config:", err)
		os.Exit(1)
	}

	s := &Core{
		config: config,
	}

	for _, option := range options {
		option(s)
	}

	return s
}

func WithConfig(config *config.Config) func(*Core) {
	return func(c *Core) {
		c.config = config
	}
}

func WithMsgChan(msgChan *chan common.Msg) func(*Core) {
	return func(c *Core) {
		c.msgChan = msgChan
	}
}

func (c *Core) getConfig() *config.Config {
	return c.config
}

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

func (c *Core) GetCommandsConfig() error {
	commands, err := c.GetCommands()

	if err != nil {
		return fmt.Errorf("error getting commands. Did you run init?")
	}

	c.commands = commands

	return nil
}

func (c *Core) GetProjectsConfig() error {
	projects, err := c.GetProjects()

	if err != nil {
		return fmt.Errorf("error getting projects. Did you run init?")
	}

	c.projects = projects

	return nil
}

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

func (c *Core) sendMsg(msg common.Msg) {
	*c.msgChan <- msg
}
