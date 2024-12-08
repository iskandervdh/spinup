package spinup

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/iskandervdh/spinup/cli"
)

type commandWithName struct {
	command string
	name    string
}

func (s *Spinup) commandTemplate(command string, project Project) string {
	// Replace placeholders in command with project values
	command = strings.ReplaceAll(command, "{{port}}", fmt.Sprintf("%d", project.Port))
	command = strings.ReplaceAll(command, "{{domain}}", project.Domain)

	for key, value := range project.Variables {
		command = strings.ReplaceAll(command, fmt.Sprintf("{{%s}}", key), value)
	}

	return command
}

func (s *Spinup) prefixOutput(prefix string, reader io.Reader, writer io.Writer) {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		fmt.Fprintf(writer, "%s %s\n", prefix, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		cli.ErrorPrint("Error reading output:", err)
	}
}

func (s *Spinup) runCommand(wg *sync.WaitGroup, project Project, command commandWithName) {
	defer wg.Done()

	cmd := exec.Command(strings.Split(command.command, " ")[0], strings.Split(command.command, " ")[1:]...)
	// Force color output
	cmd.Env = append(os.Environ(), "FORCE_COLOR=1")

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		cli.ErrorPrint("Error creating StdoutPipe:", err)
		return
	}

	stderr, err := cmd.StderrPipe()

	if err != nil {
		cli.ErrorPrint("Error creating StderrPipe:", err)
		return
	}

	go s.prefixOutput(fmt.Sprintf("[%s]", command.name), stdout, os.Stdout)
	go s.prefixOutput(fmt.Sprintf("[%s]", command.name), stderr, os.Stderr)

	// Run the project in the project's directory if it's set
	if project.Dir != nil {
		cmd.Dir = *project.Dir
	}

	err = cmd.Start()

	if err != nil {
		cli.ErrorPrint("Error starting command: ", err)
		return
	}

	err = cmd.Wait()

	if err != nil {
		// Gracefully exit if the command was interrupted by the user
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == -1 {
			return
		}

		cli.ErrorPrint("Command finished with error: ", err)
		return
	}
}

func (s *Spinup) run(project Project, projectName string) {
	var wg sync.WaitGroup
	wg.Add(len(project.Commands))

	// Start a signal listener for Ctrl+C (SIGINT)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	cli.InfoPrintf("Running project '%s'...", projectName)

	commands := []commandWithName{}

	for _, commandName := range project.Commands {
		command, err := s.getCommand(commandName)

		if err != nil {
			cli.ErrorPrintf("Error getting command '%s': %s\n", commandName, err)
			os.Exit(1)
		}

		commands = append(
			commands,
			commandWithName{
				command: s.commandTemplate(command, project),
				name:    commandName,
			})
	}

	if len(commands) == 0 {
		cli.ErrorPrint("No commands found")
		os.Exit(1)
	}

	for _, command := range commands {
		go s.runCommand(&wg, project, command)
	}

	go func() {
		<-sigChan
		cli.InfoPrintf("\nGracefully stopping project '%s'...", projectName)
	}()

	wg.Wait()
}

func (s *Spinup) tryToRun(name string) bool {
	if name == "" {
		cli.ErrorPrint("No name provided")
		os.Exit(1)
	}

	exists, project := s.projectExists(name)

	if !exists {
		return false
	}

	s.run(project, name)

	return true
}
