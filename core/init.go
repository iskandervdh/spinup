package core

import (
	"encoding/json"
	"os"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
)

// Create the config directory if it doesn't exist.
func (c *Core) createConfigDir() error {
	err := os.MkdirAll(c.config.GetConfigDir(), 0755)

	if err != nil {
		return err
	}

	return nil
}

// Create the projects.json file if it doesn't exist.
func (c *Core) createProjectsConfigFile() common.Msg {
	projectFilePath := c.getProjectsFilePath()

	if _, err := os.Stat(projectFilePath); err == nil {
		return common.NewWarnMsg("%s already exists at %s\nSkipping initialization", config.ProjectsFileName, projectFilePath)
	} else if !os.IsNotExist(err) {
		return common.NewErrMsg("Error checking if %s file exists: %v", config.ProjectsFileName, err)
	}

	emptyProjects := Projects{}
	emptyData, err := json.MarshalIndent(emptyProjects, "", "  ")

	if err != nil {
		return common.NewErrMsg("Error encoding empty projects to json: %v", err)
	}

	err = os.WriteFile(projectFilePath, emptyData, 0644)

	if err != nil {
		return common.NewErrMsg("Error writing empty projects to file: %v", err)
	}

	return common.NewInfoMsg("Initialized empty projects.json file at ", projectFilePath)
}

// Create the commands.json file if it doesn't exist.
func (c *Core) createCommandsConfigFile() common.Msg {
	commandsFilePath := c.getCommandsFilePath()

	if _, err := os.Stat(commandsFilePath); err == nil {
		return common.NewWarnMsg("%s already exists at %s\nSkipping initialization", config.CommandsFileName, commandsFilePath)
	} else if !os.IsNotExist(err) {
		return common.NewErrMsg("Error checking if commands.json file exists: %v", err)
	}

	emptyCommands := Commands{}
	emptyData, err := json.MarshalIndent(emptyCommands, "", "  ")

	if err != nil {
		return common.NewErrMsg("Error encoding empty commands to json: %v", err)
	}

	err = os.WriteFile(commandsFilePath, emptyData, 0644)

	if err != nil {
		return common.NewErrMsg("Error writing empty commands to file: %v", err)
	}

	return common.NewInfoMsg("Initialized empty commands.json file at ", commandsFilePath)
}

// TODO: Handle other types of messages and send them to the CLI via events

// Initialize the config directory and files.
func (c *Core) Init() common.Msg {
	err := c.createConfigDir()

	if err != nil {
		return common.NewErrMsg("Error creating config directory: %v", err)
	}

	msg := c.createProjectsConfigFile()

	if _, ok := msg.(*common.ErrMsg); ok {
		return common.NewErrMsg("Error creating projects.json file: %s", msg.GetText())
	} else {
		*c.msgChan <- msg
	}

	msg = c.createCommandsConfigFile()

	if _, ok := msg.(*common.ErrMsg); ok {
		return common.NewErrMsg("Error creating commands.json file: %s", msg.GetText())
	} else {
		*c.msgChan <- msg
	}

	c.RequireSudo()
	err = c.config.InitHosts()

	// if warn != nil {
	// 	c.cli.WarningPrint(warn)
	// }

	if err != nil {
		return common.NewErrMsg("Error initializing hosts file: %v", err)
	}

	err = c.config.InitNginx()

	if err != nil {
		return common.NewErrMsg("Error initializing nginx: %v", err)
	}

	return common.NewSuccessMsg("\nInitialization complete")
}
