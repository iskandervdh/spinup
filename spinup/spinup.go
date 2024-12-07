package spinup

import (
	"fmt"
	"os"
	"path"

	"github.com/iskandervdh/spinup/cli"
	"github.com/iskandervdh/spinup/config"
)

type Spinup struct {
	configDirPath string
	commands      Commands
	projects      Projects
}

func getConfigDirPath() string {
	home, err := os.UserHomeDir()

	if err != nil {
		fmt.Println("Cloud not get home directory of current user")
		panic(err)
	}

	return path.Join(home, ".config", config.ProgramName)
}

func New() *Spinup {
	return &Spinup{
		configDirPath: getConfigDirPath(),
	}
}

func (s *Spinup) getCommandsConfig() {
	commands, err := s.getCommands()

	if err != nil {
		fmt.Println("Error getting commands. Did you run init?")
		os.Exit(1)
	}

	s.commands = commands
}

func (s *Spinup) getProjectsConfig() {
	projects, err := s.getProjects()

	if err != nil {
		fmt.Println("Error getting projects. Did you run init?")
		os.Exit(1)
	}

	s.projects = projects
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
		fmt.Printf("%s %s\n", config.ProgramName, config.Version)
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
	case "cli":
		selectedOption := cli.Selection("Select option", []string{"Option 1", "Option 2", "Option 3"})
		cli.ClearTerminal()
		q := cli.Selection("Select option", []string{"Option 1", "Option 2", "Option 3"})
		cli.ClearTerminal()
		i := cli.Input("Wat wat wat?")

		fmt.Println("You selected:")
		fmt.Println(selectedOption)
		fmt.Println(q)
		fmt.Println(i)
	default:
		if !s.tryToRun(os.Args[1]) {
			fmt.Printf("Unknown subcommand or project '%s'\n", os.Args[1])
			fmt.Println("Expected 'command|c', 'project|p', 'run' or 'init' subcommand or a valid project name")
		}
	}
}
