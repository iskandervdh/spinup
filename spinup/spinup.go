package spinup

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/iskandervdh/spinup/cli"
	"github.com/iskandervdh/spinup/config"
)

type Spinup struct {
	config   config.Config
	commands Commands
	projects Projects
}

func New() *Spinup {
	return &Spinup{
		config: *config.New(),
	}
}

func (s *Spinup) requireSudo() {
	err := exec.Command("sudo", "-v").Run()

	if err != nil {
		cli.ErrorPrint("This command requires sudo")
		os.Exit(1)
	}
}

func (s *Spinup) getCommandsConfig() {
	commands, err := s.getCommands()

	if err != nil {
		cli.ErrorPrint("Error getting commands. Did you run init?")
		os.Exit(1)
	}

	s.commands = commands
}

func (s *Spinup) getProjectsConfig() {
	projects, err := s.getProjects()

	if err != nil {
		cli.ErrorPrint("Error getting projects. Did you run init?")
		os.Exit(1)
	}

	s.projects = projects
}

func (s *Spinup) getCommandNames() []string {
	var commandNames []string

	for commandName := range s.commands {
		commandNames = append(commandNames, commandName)
	}

	return commandNames
}

func (s *Spinup) getProjectNames() []string {
	var projectNames []string

	for commandName := range s.projects {
		projectNames = append(projectNames, commandName)
	}

	return projectNames
}

func (s *Spinup) Handle() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <command|project|run|init> [args...]\n", config.ProgramName)
		return
	}

	if os.Args[1] == "init" {
		s.init()
		return
	}

	s.getCommandsConfig()
	s.getProjectsConfig()

	switch os.Args[1] {
	case "-v", "--version":
		fmt.Printf("%s %s\n", config.ProgramName, strings.TrimSpace(config.Version))
	case "command", "c":
		s.handleCommand()
	case "project", "p":
		s.handleProject()
	case "variable", "v":
		s.handleVariable()
	case "run":
		if !s.tryToRun(os.Args[2]) {
			fmt.Printf("Unknown project '%s'\n", os.Args[2])
		}
	default:
		if !s.tryToRun(os.Args[1]) {
			fmt.Printf("Unknown subcommand or project '%s'\n", os.Args[1])
			fmt.Println("Expected 'command|c', 'project|p', 'run' or 'init' subcommand or a valid project name")
		}
	}
}
