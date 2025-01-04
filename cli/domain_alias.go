package cli

import (
	"fmt"
	"os"

	"github.com/iskandervdh/spinup/common"
)

// Print a list of all domain aliases for a project to the output of the CLI.
func (c *CLI) listDomainAliases(name string) error {
	exists, project := c.core.ProjectExists(name)

	if !exists {
		return fmt.Errorf("project '%s' does not exist", name)
	}

	c.sendMsg(common.NewRegularMsg("%-20s %-30s\n", "Key", "Value"))

	for key, value := range project.Variables {
		c.sendMsg(common.NewRegularMsg("%-20s %-30s\n", key, value))
	}

	return nil
}

// Handle the domain-alias command.
func (c *CLI) handleDomainAlias() {
	if len(os.Args) < 3 {
		c.sendMsg(common.NewRegularMsg("Usage: %s domain-alias|da <add|remove|list> [args...]\n", common.ProgramName))
		return
	}

	switch os.Args[2] {
	case "list", "ls":
		if len(os.Args) < 4 {
			c.sendMsg(common.NewRegularMsg("Usage: %s domain-alias|da list|ls <project>\n", common.ProgramName))
			return
		}

		err := c.listDomainAliases(os.Args[3])

		if err != nil {
			c.ErrorPrint("Error listing variables:", err)
		}
	case "add":
		if len(os.Args) < 5 {
			c.sendMsg(common.NewRegularMsg("Usage: %s domain-alias|da add <project> <domain-alias>\n", common.ProgramName))
			return
		}

		c.sendMsg(c.core.AddDomainAlias(os.Args[3], os.Args[4]))
	case "remove", "rm":
		if len(os.Args) < 5 {
			c.sendMsg(common.NewRegularMsg("Usage: %s domain-alias remove|rm <project> <domain-alias>\n", common.ProgramName))
			return
		}

		c.sendMsg(c.core.RemoveDomainAlias(os.Args[3], os.Args[4]))
	}
}
