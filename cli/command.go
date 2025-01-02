package cli

import (
	"fmt"
	"os"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
)

// Print a list of all commands to the output of the CLI.
func (c *CLI) listCommands() {
	commands, err := c.core.GetCommands()

	if err != nil {
		c.ErrorPrint("Error getting commands:", err)
		return
	}

	fmt.Fprintf(c.out, "%-20s %-30s\n", "Name", "Command")

	for _, command := range commands {
		fmt.Fprintf(c.out, "%-20s %-30s\n", command.Name, command.Command)
	}
}

// Add a command interactively by asking the user for the name and command.
func (c *CLI) addCommandInteractive() {
	name := c.Input("Enter command name:", "")
	command := c.Input("Enter command:", "")

	c.sendMsg(c.core.AddCommand(name, command))
}

// Remove a command interactively by asking the user to select a command to remove.
func (c *CLI) removeCommandInteractive() {
	name, err, exited := c.Selection("Select command to remove", c.core.GetCommandNames())

	if err != nil {
		c.ErrorPrint("Error selecting command:", err)
		return
	}

	if exited {
		return
	}

	if name == "" {
		c.ErrorPrint("No command selected")
		return
	}

	if !c.Confirm("Are you sure you want to remove command " + name + "?") {
		return
	}

	c.sendMsg(c.core.RemoveCommand(name))
}

// Edit a command interactively by asking the user to select a command to edit and then enter a new command.
func (c *CLI) editCommandInteractive() {
	name, err, exited := c.Selection("Select command to edit", c.core.GetCommandNames())

	if err != nil {
		c.ErrorPrint("Error selecting command:", err)
		return
	}

	if exited {
		return
	}

	if name == "" {
		c.ErrorPrint("No command selected")
		return
	}

	exist, command := c.core.CommandExists(name)

	if !exist {
		c.ErrorPrint("Command does not exist")
		return
	}

	newCommand := c.Input("Edit command:", command.Command)

	if !c.Confirm("Are you sure you want to update command " + name + "?") {
		return
	}

	c.sendMsg(c.core.UpdateCommand(name, newCommand))
}

// Handle the command subcommand.
func (c *CLI) handleCommand() {
	if len(os.Args) < 3 {
		c.sendMsg(common.NewRegularMsg("Usage: spinup command <add|remove|edit|rename|list> [args...]\n"))
		return
	}

	commandName := os.Args[2]

	switch commandName {
	case "list", "ls":
		c.listCommands()
	case "add":
		if len(os.Args) == 3 {
			c.addCommandInteractive()
			return
		}

		if len(os.Args) < 5 {
			c.sendMsg(common.NewRegularMsg("Usage: %s command|c add <name> <command>\n", config.ProgramName))
			return
		}

		c.sendMsg(c.core.AddCommand(os.Args[3], os.Args[4]))
	case "remove", "rm":
		if len(os.Args) == 3 {
			c.removeCommandInteractive()
			return
		}

		if len(os.Args) < 4 {
			c.sendMsg(common.NewRegularMsg("Usage: %s command|c remove|rm <name>\n", config.ProgramName))
			return
		}

		c.sendMsg(c.core.RemoveCommand(os.Args[3]))
	case "edit", "e":
		if len(os.Args) == 3 {
			c.editCommandInteractive()
			return
		}

		if len(os.Args) < 5 {
			c.sendMsg(common.NewRegularMsg("Usage: %s command|c edit|e <name> <command>\n", config.ProgramName))
			return
		}

		c.sendMsg(c.core.UpdateCommand(os.Args[3], os.Args[4]))
	case "rename", "mv":
		if len(os.Args) < 5 {
			c.sendMsg(common.NewRegularMsg("Usage: %s command|c rename|mv <old-name> <new-name>\n", config.ProgramName))
			return
		}

		c.sendMsg(c.core.RenameCommand(os.Args[3], os.Args[4]))
	default:
		c.sendMsg(common.NewErrMsg("Unknown subcommand '%s'\n", commandName))
		c.sendMsg(common.NewRegularMsg("Expected 'add', 'remove', 'edit', 'rename' or 'list'\n"))
	}
}
