package core

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

	"github.com/iskandervdh/spinup/common"
)

// commandWithName is a struct to hold a command and its name.
type commandWithName struct {
	command string
	name    string
}

func (c *Core) commandTemplate(command string, project Project) string {
	// Replace placeholders in command with project values
	command = strings.ReplaceAll(command, "{{port}}", fmt.Sprintf("%d", project.Port))
	command = strings.ReplaceAll(command, "{{domain}}", project.Domain)

	for key, value := range project.Variables {
		command = strings.ReplaceAll(command, fmt.Sprintf("{{%s}}", key), value)
	}

	return command
}

func (c *Core) prefixOutput(prefix string, reader io.Reader, writer io.Writer) error {
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		fmt.Fprintf(writer, "%s %s\n", prefix, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading output: %s", err)
	}

	return nil
}

func (c *Core) runCommand(wg *sync.WaitGroup, project Project, command commandWithName) error {
	defer wg.Done()

	cmd := exec.Command(strings.Split(command.command, " ")[0], strings.Split(command.command, " ")[1:]...)
	// Force color output
	cmd.Env = append(os.Environ(), "FORCE_COLOR=1")

	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return fmt.Errorf("error creating StdoutPipe: %s", err)
	}

	stderr, err := cmd.StderrPipe()

	if err != nil {
		return fmt.Errorf("error creating StderrPipe: %s", err)
	}

	go c.prefixOutput(fmt.Sprintf("[%s]", command.name), stdout, os.Stdout)
	go c.prefixOutput(fmt.Sprintf("[%s]", command.name), stderr, os.Stderr)

	// Run the project in the project's directory if it's set
	if project.Dir != nil {
		cmd.Dir = *project.Dir
	}

	err = cmd.Start()

	if err != nil {
		return fmt.Errorf("error starting command: %s", err)
	}

	err = cmd.Wait()

	if err != nil {
		// Gracefully exit if the command was interrupted by the user
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == -1 {
			return nil
		}

		return fmt.Errorf("command finished with error: %s", err)
	}

	return nil
}

// Run a project with the given name.
func (c *Core) run(project Project, projectName string) common.Msg {
	var wg sync.WaitGroup
	wg.Add(len(project.Commands))

	// Start a signal listener for Ctrl+C (SIGINT) to gracefully stop the project when the user interrupts the process.
	c.sigChan = make(chan os.Signal, 1)
	signal.Notify(c.sigChan, syscall.SIGINT, syscall.SIGTERM)

	c.sendMsg(common.NewInfoMsg("Running project '%s'...", projectName))

	commands := []commandWithName{}

	// Add all commands to the commands array in a form that includes the command name.
	for _, commandName := range project.Commands {
		exists, command := c.CommandExists(commandName)

		if !exists {
			return common.NewErrMsg("Command '%s' does not exist", commandName)
		}

		commands = append(
			commands,
			commandWithName{
				command: c.commandTemplate(command, project),
				name:    commandName,
			})
	}

	if len(commands) == 0 {
		return common.NewErrMsg("No commands found")
	}

	for _, command := range commands {
		go c.runCommand(&wg, project, command)
	}

	go func() {
		<-c.sigChan

		c.sendMsg(common.NewInfoMsg("\nGracefully stopping project '%s'...", projectName))
		wg.Done()
	}()

	wg.Wait()

	return common.NewSuccessMsg("")
}

// Try to run a project with the given name.
func (c *Core) TryToRun(name string) common.Msg {
	if name == "" {
		return common.NewErrMsg("No name provided")
	}

	exists, project := c.ProjectExists(name)

	if !exists {
		return nil
	}

	return c.run(project, name)
}
