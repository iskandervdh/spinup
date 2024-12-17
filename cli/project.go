package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
)

func (c *CLI) listProjects() {
	projects, err := c.core.GetProjects()

	if err != nil {
		c.ErrorPrint("Error getting commands:", err)
		return
	}

	c.sendMsg(common.NewRegularMsg("%-10s %-30s %-10s %-20s\n", "Name", "Domain", "Port", "Commands"))

	for projectName, project := range projects {
		c.sendMsg(
			common.NewRegularMsg("%-10s %-30s %-10d %-20s\n",
				projectName,
				project.Domain,
				project.Port,
				strings.Join(project.Commands, ", "),
			),
		)
	}
}

func (c *CLI) addProject(name string, domain string, port int, commandNames []string) {
	c.Loading(fmt.Sprintf("Adding project %s...", name),
		func() common.Msg {
			return c.core.AddProject(name, domain, port, commandNames)
		},
	)
}

func (c *CLI) addProjectInteractive() {
	name := c.Input("Project name:")
	domain := c.Input("Domain:")
	port := c.Input("Port:")

	portInt, err := strconv.Atoi(port)

	if err != nil {
		c.ErrorPrint("Port must be an integer")
		return
	}

	selectedCommands, err, exited := c.Question("Commands", c.core.GetCommandNames())

	if err != nil {
		c.ErrorPrint("Error getting commands:", err)
		os.Exit(1) // TODO: Maybe return error instead of exiting
	}

	if exited {
		return
	}

	c.addProject(name, domain, portInt, selectedCommands)
}

func (c *CLI) RemoveProject(name string) {
	c.Loading(fmt.Sprintf("Removing project %s...", name),
		func() common.Msg {
			return c.core.RemoveProject(name)
		},
	)
}

func (c *CLI) removeProjectInteractive() {
	name, err, exited := c.Selection("What project do you want to remove?", c.core.GetProjectNames())

	if err != nil {
		c.ErrorPrint("Error getting project names:", err)
		return
	}

	if exited {
		return
	}

	if name == "" {
		c.ErrorPrint("No project selected")
		return
	}

	if !c.Confirm("Are you sure you want to remove project " + name + "?") {
		return
	}

	c.core.RemoveProject(name)
}

func (c *CLI) handleProject() {
	if len(os.Args) < 3 {
		c.sendMsg(common.NewInfoMsg("Usage: %s project <add|remove|list> [args...]", config.ProgramName))
		return
	}

	var msg common.Msg

	switch os.Args[2] {
	case "list", "ls":
		c.listProjects()
	case "add":
		if len(os.Args) == 3 {
			c.addProjectInteractive()
			return
		}

		if len(os.Args) < 6 {
			c.sendMsg(common.NewRegularMsg("Usage: %s project add <name> <domain> <port>\n", config.ProgramName))
			return
		}

		port, err := strconv.Atoi(os.Args[5])

		if err != nil {
			c.ErrorPrint("Port must be an integer")
			return
		}

		c.addProject(os.Args[3], os.Args[4], port, os.Args[6:])
	case "remove", "rm":
		if len(os.Args) == 3 {
			c.removeProjectInteractive()
			return
		}

		if len(os.Args) != 4 {
			c.sendMsg(common.NewRegularMsg("Usage: %s project remove|rm <name>\n", config.ProgramName))
			return
		}

		c.RemoveProject(os.Args[3])
	case "add-command", "ac":
		if len(os.Args) < 5 {
			c.sendMsg(common.NewRegularMsg("Usage: %s project add-command|ac <project> <command>\n", config.ProgramName))
			return
		}

		msg = c.core.AddCommandToProject(os.Args[3], os.Args[4])
	case "remove-command", "rc":
		if len(os.Args) < 5 {
			c.sendMsg(common.NewRegularMsg("Usage: %s project remove-command|rc <project> <command>\n", config.ProgramName))
			return
		}

		msg = c.core.RemoveCommandFromProject(os.Args[3], os.Args[4])
	case "set-dir", "sd":
		if len(os.Args) < 4 {
			c.sendMsg(common.NewRegularMsg("Usage: %s project set-dir|sp <project> [dir]\n", config.ProgramName))
			return
		}

		if len(os.Args) == 5 {
			c.core.SetProjectDir(os.Args[3], &os.Args[4])
			return
		}

		msg = c.core.SetProjectDir(os.Args[3], nil)
	case "get-dir", "gd":
		if len(os.Args) != 4 {
			c.sendMsg(common.NewRegularMsg("Usage: %s project get-dir|gp <project>\n", config.ProgramName))
			return
		}

		msg = c.core.GetProjectDir(os.Args[3])
	default:
		msg = common.NewErrMsg("Unknown subcommand '%s'", os.Args[2])

	}

	c.MsgPrint(msg)
}
