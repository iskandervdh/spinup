package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
)

type Commands map[string]string

func (c *Core) getCommandsFilePath() string {
	return path.Join(c.config.GetConfigDir(), config.CommandsFileName)
}

func (c *Core) GetCommands() (Commands, error) {
	commandsFileContent, err := os.ReadFile(c.getCommandsFilePath())

	if err != nil {
		return nil, fmt.Errorf("error reading commands.json file: %s", err)
	}

	var commands Commands
	err = json.Unmarshal(commandsFileContent, &commands)

	if err != nil {
		return nil, fmt.Errorf("error parsing commands.json file: %s", err)
	}

	return commands, nil
}

func (c *Core) getCommand(name string) (string, error) {
	command, exists := c.commands[name]

	if !exists {
		return "", fmt.Errorf("command '%s' not found", name)
	}

	return command, nil
}

func (c *Core) AddCommand(name string, command string) common.Msg {
	// Check if already exists
	for commandName := range c.commands {
		if commandName == name {
			return common.NewErrMsg("command '%s' already exists", name)
		}
	}

	c.commands[name] = command

	updatedCommands, err := json.MarshalIndent(c.commands, "", "  ")

	if err != nil {
		return common.NewErrMsg("error encoding commands to json: %s", err)
	}

	err = os.WriteFile(c.getCommandsFilePath(), updatedCommands, 0644)

	if err != nil {
		return common.NewErrMsg("error writing commands to file: %s", err)
	}

	return common.NewSuccessMsg("Added command '%s': %s", name, command)
}

func (c *Core) RemoveCommand(name string) common.Msg {
	if c.commands == nil {
		return common.NewErrMsg("No commands found")
	}

	delete(c.commands, name)

	updatedCommands, err := json.MarshalIndent(c.commands, "", "  ")

	if err != nil {
		return common.NewErrMsg("Error encoding commands to json: %s", err)
	}

	err = os.WriteFile(c.getCommandsFilePath(), updatedCommands, 0644)

	if err != nil {
		return common.NewErrMsg("Error writing commands to file: %s", err)
	}

	return common.NewSuccessMsg("Removed command '%s'", name)
}
