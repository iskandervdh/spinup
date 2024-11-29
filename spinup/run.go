package spinup

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (s *Spinup) commandTemplate(command string, project Project) string {
	// Replace placeholders in command with project values
	command = strings.ReplaceAll(command, "{{port}}", fmt.Sprintf("%d", project.Port))
	command = strings.ReplaceAll(command, "{{domain}}", project.Domain)

	for key, value := range project.Variables {
		command = strings.ReplaceAll(command, fmt.Sprintf("{{%s}}", key), value)
	}

	return command
}

func (s *Spinup) run(project Project, projectName string) {
	fmt.Printf("Running project %s\n", projectName)

	commands := []string{}

	for _, commandName := range project.Commands {
		command, err := s.getCommand(commandName)

		if err != nil {
			fmt.Printf("Error getting command '%s': %s\n", commandName, err)
			os.Exit(1)
		}

		commands = append(commands, s.commandTemplate(command, project))
	}

	if len(commands) == 0 {
		fmt.Println("No commands found")
		os.Exit(1)
	}

	concurrent := exec.Command("concurrently", commands...)
	concurrent.Stdout = os.Stdout

	concurrent.Run()
}

func (s *Spinup) tryToRun(name string) bool {
	if name == "" {
		fmt.Println("No name provided")
		os.Exit(1)
	}

	exists, project := s.projectExists(name)

	if !exists {
		return false
	}

	s.run(project, name)

	return true
}
