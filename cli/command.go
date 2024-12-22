package cli

import (
	"fmt"
	"os"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
)

func (c *CLI) listCommands() {
	commands, err := c.core.GetCommands()

	if err != nil {
		c.ErrorPrint("Error getting commands:", err)
		return
	}

	fmt.Fprintf(c.out, "%-20s %-30s\n", "Name", "Command")

	for commandName, command := range commands {
		fmt.Fprintf(c.out, "%-20s %-30s\n", commandName, command)
	}
}

func (c *CLI) addCommandInteractive() {
	name := c.Input("Enter command name:", "")
	command := c.Input("Enter command:", "")

	c.MsgPrint(c.core.AddCommand(name, command))
}

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

	c.MsgPrint(c.core.RemoveCommand(name))
}

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

	newCommand := c.Input("Edit command:", command)

	if !c.Confirm("Are you sure you want to update command " + name + "?") {
		return
	}

	c.MsgPrint(c.core.UpdateCommand(name, newCommand))
}

func (c *CLI) handleCommand() {
	if len(os.Args) < 3 {
		c.sendMsg(common.NewInfoMsg("Usage: spinup command <add|remove|list> [args...]"))
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
			c.sendMsg(common.NewInfoMsg("Usage: %s command add <name> <command>", config.ProgramName))
			return
		}

		c.MsgPrint(c.core.AddCommand(os.Args[3], os.Args[4]))
	case "remove", "rm":
		if len(os.Args) == 3 {
			c.removeCommandInteractive()
			return
		}

		if len(os.Args) < 4 {
			c.sendMsg(common.NewInfoMsg("Usage: %s command remove <name>", config.ProgramName))
			return
		}

		c.MsgPrint(c.core.RemoveCommand(os.Args[3]))
	case "edit", "e":
		if len(os.Args) == 3 {
			c.editCommandInteractive()
			return
		}

		if len(os.Args) < 5 {
			c.sendMsg(common.NewInfoMsg("Usage: %s command edit <name> <command>", config.ProgramName))
			return
		}

		c.MsgPrint(c.core.UpdateCommand(os.Args[3], os.Args[4]))
	case "rename", "mv":
		if len(os.Args) < 5 {
			c.sendMsg(common.NewInfoMsg("Usage: %s command rename <old-name> <new-name>", config.ProgramName))
			return
		}

		c.MsgPrint(c.core.RenameCommand(os.Args[3], os.Args[4]))
	default:
		c.sendMsg(common.NewErrMsg("Unknown subcommand '%s'\n", commandName))
		c.sendMsg(common.NewInfoMsg("Expected 'add', 'remove' or 'list'"))
	}
}
