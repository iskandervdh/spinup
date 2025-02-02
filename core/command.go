package core

import (
	"fmt"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/database/sqlc"
)

type Command = sqlc.Command

type Commands []Command

func (c *Core) FetchCommands() error {
	commands, err := c.dbQueries.GetCommands(c.dbContext)

	if err != nil {
		return fmt.Errorf("error getting commands: %s", err)
	}

	c.commands = commands

	return nil
}

// Get the commands from the database.
func (c *Core) GetCommands() ([]Command, error) {
	if c.commands == nil {
		err := c.FetchCommands()

		if err != nil {
			return nil, err
		}
	}

	return c.commands, nil
}

// Check if a command with the given name exists. Returns the command if it exists.
func (c *Core) CommandExists(name string) (bool, Command) {
	if c.commands == nil {
		err := c.FetchCommands()

		if err != nil {
			return false, Command{}
		}
	}

	command, err := c.dbQueries.GetCommand(c.dbContext, name)

	if err != nil {
		return false, Command{}
	}

	return true, command
}

// Add a command with the given name and command string.
func (c *Core) AddCommand(name string, command string) common.Msg {
	// Check if already exists
	for _, command := range c.commands {
		if command.Name == name {
			return common.NewErrMsg("command '%s' already exists", name)
		}
	}

	err := c.dbQueries.CreateCommand(c.dbContext, sqlc.CreateCommandParams{
		Name:    name,
		Command: command,
	})

	if err != nil {
		return common.NewErrMsg("error adding command to database: %s", err)
	}

	return common.NewSuccessMsg("Added command '%s': %s", name, command)
}

// Remove the command with the given name.
func (c *Core) RemoveCommand(name string) common.Msg {
	if c.commands == nil {
		return common.NewErrMsg("No commands found")
	}

	err := c.dbQueries.DeleteCommand(c.dbContext, name)

	if err != nil {
		return common.NewErrMsg("Error deleting command: %s", err)
	}

	return common.NewSuccessMsg("Removed command '%s'", name)
}

func (c *Core) RemoveCommandById(id int64) common.Msg {
	if c.commands == nil {
		return common.NewErrMsg("No commands found")
	}

	err := c.dbQueries.DeleteCommandById(c.dbContext, id)

	if err != nil {
		return common.NewErrMsg("Error deleting command: %s", err)
	}

	return common.NewSuccessMsg("Removed command")
}

// Update the command with the given name to the given command string.
func (c *Core) UpdateCommand(name string, command string) common.Msg {
	if c.commands == nil {
		return common.NewErrMsg("No commands found")
	}

	err := c.dbQueries.UpdateCommand(c.dbContext, sqlc.UpdateCommandParams{
		Name:    name,
		Command: command,
	})

	if err != nil {
		return common.NewErrMsg("Error updating command: %s", err)
	}

	return common.NewSuccessMsg("Updated command '%s': %s", name, command)
}

func (c *Core) UpdateCommandById(id int64, name string, command string) common.Msg {
	err := c.dbQueries.UpdateCommandById(c.dbContext, sqlc.UpdateCommandByIdParams{
		ID:      id,
		Name:    name,
		Command: command,
	})

	if err != nil {
		return common.NewErrMsg("Error updating command: %s", err)
	}

	return common.NewSuccessMsg("Updated command '%s': %s", name, command)
}

// Rename the command with the given old name to the given new name.
func (c *Core) RenameCommand(oldName string, newName string) common.Msg {
	if c.commands == nil {
		return common.NewErrMsg("No commands found")
	}

	err := c.dbQueries.RenameCommand(c.dbContext, sqlc.RenameCommandParams{
		Name:   newName,
		Name_2: oldName,
	})

	if err != nil {
		return common.NewErrMsg("Error renaming command: %s", err)
	}

	return common.NewSuccessMsg("Renamed command '%s' to '%s'", oldName, newName)
}
