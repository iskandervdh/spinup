package cli

import (
	"fmt"
	"os"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
)

func (c *CLI) listVariables(name string) error {
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

func (c *CLI) handleVariable() {
	if len(os.Args) < 3 {
		c.sendMsg(common.NewRegularMsg("Usage: %s variable <add|remove|list> [args...]\n", config.ProgramName))
		return
	}

	switch os.Args[2] {
	case "list", "ls":
		if len(os.Args) < 4 {
			c.sendMsg(common.NewRegularMsg("Usage: %s variable list|ls <project>\n", config.ProgramName))
			return
		}

		err := c.listVariables(os.Args[3])

		if err != nil {
			c.ErrorPrint("Error listing variables:", err)
		}
	case "add":
		if len(os.Args) < 6 {
			c.sendMsg(common.NewRegularMsg("Usage: %s variable add <project> <key> <value>\n", config.ProgramName))
			return
		}

		c.core.AddVariable(os.Args[3], os.Args[4], os.Args[5])
	case "remove", "rm":
		if len(os.Args) < 5 {
			c.sendMsg(common.NewRegularMsg("Usage: %s variable remove|rm <project> <key>\n", config.ProgramName))
			return
		}

		c.sendMsg(c.core.RemoveVariable(os.Args[3], os.Args[4]))
	}
}
