package cli

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
)

// Print a list of all projects to the output of the CLI.
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

// Add a project and display a loading message.
func (c *CLI) addProject(name string, domain string, port int, commandNames []string) {
	c.Loading(fmt.Sprintf("Adding project %s...", name),
		func() common.Msg {
			return c.core.AddProject(name, domain, port, commandNames)
		},
	)
}

// Add a project interactively by asking the user for the name, domain, port and commands.
func (c *CLI) addProjectInteractive() {
	name := c.Input("Project name:", "")
	domain := c.Input("Domain:", "")
	port := c.Input("Port:", "")

	portInt, err := strconv.Atoi(port)

	if err != nil {
		c.ErrorPrint("Port must be an integer")
		return
	}

	selectedCommands, err, exited := c.Question("Commands", c.core.GetCommandNames(), nil)

	if err != nil {
		c.ErrorPrint("Error getting commands:", err)
		return
	}

	if exited {
		return
	}

	c.addProject(name, domain, portInt, selectedCommands)
}

// Remove a project and display a loading message.
func (c *CLI) removeProject(name string) {
	c.Loading(fmt.Sprintf("Removing project %s...", name),
		func() common.Msg {
			return c.core.RemoveProject(name)
		},
	)
}

// Remove a project interactively by asking the user to select a project to remove.
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

// Edit a project and display a loading message.
func (c *CLI) editProject(name string, domain string, port int, commandNames []string) {
	c.Loading(fmt.Sprintf("Updating project %s...", name),
		func() common.Msg {
			return c.core.UpdateProject(name, domain, port, commandNames)
		},
	)
}

// Edit a project interactively by asking the user to select a project to edit and then enter new values.
func (c *CLI) editProjectInteractive() {
	name, err, exited := c.Selection("What project do you want to edit?", c.core.GetProjectNames())

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

	exists, project := c.core.ProjectExists(name)

	if !exists {
		c.ErrorPrint("Project does not exist")
		return
	}

	domain := c.Input("Domain:", project.Domain)
	port := c.Input("Port:", strconv.Itoa(project.Port))

	portInt, err := strconv.Atoi(port)

	if err != nil {
		c.ErrorPrint("Port must be an integer")
		return
	}

	projectSelectedCommands := make([]bool, len(c.core.GetCommandNames()))
	commandNames := c.core.GetCommandNames()

	for i, commandName := range commandNames {
		projectSelectedCommands[i] = slices.Contains(project.Commands, commandName)
	}

	selectedCommands, err, exited := c.Question("Commands", c.core.GetCommandNames(), projectSelectedCommands)

	if err != nil {
		c.ErrorPrint("Error getting commands:", err)
		return
	}

	if exited {
		return
	}

	c.sendMsg(c.core.UpdateProject(name, domain, portInt, selectedCommands))
}

// Handle the project subcommand.
func (c *CLI) handleProject() {
	if len(os.Args) < 3 {
		c.sendMsg(common.NewRegularMsg("Usage: %s project <add|remove|edit|rename|add-command|remove-command|set-dir|get-dir|list> [args...]\n", config.ProgramName))
		return
	}

	switch os.Args[2] {
	case "list", "ls":
		c.listProjects()
	case "add":
		if len(os.Args) == 3 {
			c.addProjectInteractive()
			return
		}

		if len(os.Args) < 6 {
			c.sendMsg(common.NewRegularMsg("Usage: %s project add <name> <domain> <port> [command names...]\n", config.ProgramName))
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

		c.removeProject(os.Args[3])
	case "edit", "e":
		if len(os.Args) < 4 {
			c.editProjectInteractive()
			return
		}

		if len(os.Args) < 6 {
			c.sendMsg(common.NewRegularMsg("Usage: %s project edit <name> <domain> <port> [command names...]\n", config.ProgramName))
			return
		}

		port, err := strconv.Atoi(os.Args[5])

		if err != nil {
			c.ErrorPrint("Port must be an integer")
			return
		}

		c.editProject(os.Args[3], os.Args[4], port, os.Args[6:])
	case "rename", "mv":
		if len(os.Args) < 5 {
			c.sendMsg(common.NewRegularMsg("Usage: %s project rename|mv <old-name> <new-name>\n", config.ProgramName))
			return
		}

		c.sendMsg(c.core.RenameProject(os.Args[3], os.Args[4]))
	case "add-command", "ac":
		if len(os.Args) < 5 {
			c.sendMsg(common.NewRegularMsg("Usage: %s project add-command|ac <project> <command>\n", config.ProgramName))
			return
		}

		c.sendMsg(c.core.AddCommandToProject(os.Args[3], os.Args[4]))
	case "remove-command", "rc":
		if len(os.Args) < 5 {
			c.sendMsg(common.NewRegularMsg("Usage: %s project remove-command|rc <project> <command>\n", config.ProgramName))
			return
		}

		c.sendMsg(c.core.RemoveCommandFromProject(os.Args[3], os.Args[4]))
	case "set-dir", "sd":
		if len(os.Args) < 4 {
			c.sendMsg(common.NewRegularMsg("Usage: %s project set-dir|sd <project> [dir]\n", config.ProgramName))
			return
		}

		if len(os.Args) == 5 {
			c.sendMsg(c.core.SetProjectDir(os.Args[3], &os.Args[4]))
			return
		}

		c.sendMsg(c.core.SetProjectDir(os.Args[3], nil))
	case "get-dir", "gd":
		if len(os.Args) != 4 {
			c.sendMsg(common.NewRegularMsg("Usage: %s project get-dir|gd <project>\n", config.ProgramName))
			return
		}

		c.sendMsg(c.core.GetProjectDir(os.Args[3]))
	default:
		c.sendMsg(common.NewErrMsg("Unknown subcommand '%s'", os.Args[2]))
		c.sendMsg(common.NewRegularMsg("Expected 'add', 'remove|rm', 'edit|e', 'rename|mv', 'add-command|ac', 'remove-command|rc', 'set-dir|sd', 'get-dir|gd' subcommand\n"))
	}
}
