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

type runningCommand struct {
	command string
	name    string
	cmd     *exec.Cmd
}

func (c *Core) commandTemplate(command string, project Project) string {
	// Replace placeholders in command with project values
	command = strings.ReplaceAll(command, "{{port}}", fmt.Sprintf("%d", project.Port))
	command = strings.ReplaceAll(command, "{{domain}}", common.GetDomain(project.Name))

	for _, variable := range project.Variables {
		command = strings.ReplaceAll(command, fmt.Sprintf("{{%s}}", variable.Name), variable.Value)
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

func (c *Core) runCommand(wg *sync.WaitGroup, project Project, command *runningCommand) error {
	defer wg.Done()

	command.cmd = exec.Command(strings.Split(command.command, " ")[0], strings.Split(command.command, " ")[1:]...)

	// create a new process group for the command
	command.cmd.SysProcAttr = createProcessGroup()

	// Force color output
	command.cmd.Env = append(os.Environ(), "FORCE_COLOR=1")

	stdout, err := command.cmd.StdoutPipe()

	if err != nil {
		return fmt.Errorf("error creating StdoutPipe: %s", err)
	}

	stderr, err := command.cmd.StderrPipe()

	if err != nil {
		return fmt.Errorf("error creating StderrPipe: %s", err)
	}

	go c.prefixOutput(fmt.Sprintf("[%s]", command.name), stdout, c.out)
	go c.prefixOutput(fmt.Sprintf("[%s]", command.name), stderr, c.err)

	// Run the project in the project's directory if it's set
	if project.Dir.Valid {
		command.cmd.Dir = project.Dir.String
	}

	err = command.cmd.Start()

	if err != nil {
		return fmt.Errorf("error starting command: %s", err)
	}

	err = command.cmd.Wait()

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
	sigChan := make(chan os.Signal, 1)
	c.sigChan = &sigChan
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	c.sendMsg(common.NewInfoMsg("Running project '%s'...", projectName))

	runningCommands := []*runningCommand{}

	// Add all commands to the commands array in a form that includes the command name.
	for _, command := range project.Commands {
		runningCommands = append(
			runningCommands,
			&runningCommand{
				command: c.commandTemplate(command.Command, project),
				name:    command.Name,
			})
	}

	if len(runningCommands) == 0 {
		return common.NewErrMsg("No commands found")
	}

	for _, runningCommand := range runningCommands {
		go c.runCommand(&wg, project, runningCommand)
	}

	go func() {
		<-*c.sigChan

		c.sendMsg(common.NewInfoMsg("\nGracefully stopping project '%s'...", projectName))

		// Send terminate signal to all running commands
		for _, runningCommand := range runningCommands {
			if runningCommand.cmd.Process != nil {
				err := killProcess(runningCommand.cmd.Process)

				if err != nil {
					c.sendMsg(common.NewErrMsg("Failed to send SIGTERM to command '%s': %s", runningCommand.name, err))
				}
			}
		}
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
