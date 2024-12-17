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
	config   *config.Config
	cli      *cli.CLI
	commands Commands
	projects Projects
}

func New(options ...func(*Spinup)) *Spinup {
	config, err := config.New()
	cli := cli.New()

	if err != nil {
		fmt.Println("Error getting config:", err)
		os.Exit(1)
	}

	s := &Spinup{
		config: config,
		cli:    cli,
	}

	for _, option := range options {
		option(s)
	}

	return s
}

func WithConfig(cfg *config.Config) func(*Spinup) {
	return func(s *Spinup) {
		s.config = cfg
	}
}

func WithCLI(cli *cli.CLI) func(*Spinup) {
	return func(s *Spinup) {
		s.cli = cli
	}
}

func (s *Spinup) getConfig() *config.Config {
	return s.config
}

func (s *Spinup) requireSudo() {
	if s.config.IsTesting() {
		return
	}

	err := exec.Command("sudo", "-v").Run()

	if err != nil {
		s.cli.ErrorPrint("This command requires sudo")
		os.Exit(1)
	}
}

func (s *Spinup) getCommandsConfig() {
	commands, err := s.getCommands()

	if err != nil {
		s.cli.ErrorPrint("Error getting commands. Did you run init?")
		os.Exit(1)
	}

	s.commands = commands
}

func (s *Spinup) getProjectsConfig() {
	projects, err := s.getProjects()

	if err != nil {
		s.cli.ErrorPrint("Error getting projects. Did you run init?")
		os.Exit(1)
	}

	s.projects = projects
}

func (s *Spinup) getCommandNames() []string {
	if s.commands == nil {
		s.getCommandsConfig()
	}

	var commandNames []string

	for commandName := range s.commands {
		commandNames = append(commandNames, commandName)
	}

	return commandNames
}

func (s *Spinup) getProjectNames() []string {
	if s.projects == nil {
		s.getProjectsConfig()
	}

	var projectNames []string

	for commandName := range s.projects {
		projectNames = append(projectNames, commandName)
	}

	return projectNames
}

func (s *Spinup) getCommandsForProject(projectName string) []string {
	if s.projects == nil {
		s.getProjectsConfig()
	}

	project, ok := s.projects[projectName]

	if !ok {
		return nil
	}

	return project.Commands
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
