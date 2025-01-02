package core

import (
	"sort"
	"testing"
)

func TestGetCommandNames(t *testing.T) {
	c := TestingCore("get_command_names")

	commandNames := c.GetCommandNames()

	if len(commandNames) != 0 {
		t.Error("Expected no command names, got", len(commandNames))
		return
	}
}

func TestAddCommand(t *testing.T) {
	c := TestingCore("add_command")

	err := c.FetchCommands()

	if err != nil {
		t.Error("Error fetching commands:", err)
		return
	}

	c.AddCommand("test", "ls")

	// "Refetch" the commands config
	err = c.FetchCommands()

	if err != nil {
		t.Error("Error fetching commands:", err)
		return
	}

	commandNames := c.GetCommandNames()

	if len(commandNames) != 1 {
		t.Error("Expected 1 command name, got", len(commandNames))
		return
	}

	if commandNames[0] != "test" {
		t.Error("Expected command name to be 'test', got", commandNames[0])
		return
	}
}

func TestAddCommandDuplicate(t *testing.T) {
	c := TestingCore("add_command_duplicate")

	c.FetchCommands()

	c.AddCommand("test", "ls")

	// "Refetch" the commands config
	c.FetchCommands()

	c.AddCommand("test", "ls")

	// "Refetch" the commands config
	c.FetchCommands()

	commandNames := c.GetCommandNames()

	if len(commandNames) != 1 {
		t.Error("Expected 1 command name, got", len(commandNames))
		return
	}

	if commandNames[0] != "test" {
		t.Error("Expected command name to be 'test', got", commandNames[0])
		return
	}
}

func TestAddMultipleCommands(t *testing.T) {
	c := TestingCore("add_multiple_commands")

	c.FetchCommands()

	c.AddCommand("test", "ls")
	c.AddCommand("test2", "ls")

	// "Refetch" the commands config
	c.FetchCommands()

	commandNames := c.GetCommandNames()
	sort.Strings(commandNames)

	if len(commandNames) != 2 {
		t.Error("Expected 2 command names, got", len(commandNames))
		return
	}

	if commandNames[0] != "test" {
		t.Error("Expected command name to be 'test', got", commandNames[0])
		return
	}

	if commandNames[1] != "test2" {
		t.Error("Expected command name to be 'test2', got", commandNames[1])
		return
	}
}

func TestRemoveCommand(t *testing.T) {
	c := TestingCore("remove_command")

	c.FetchCommands()

	c.AddCommand("test", "ls")

	// "Refetch" the commands config
	c.FetchCommands()

	c.RemoveCommand("test")

	// "Refetch" the commands config
	c.FetchCommands()

	commandNames := c.GetCommandNames()

	if len(commandNames) != 0 {
		t.Error("Expected no command names, got", len(commandNames))
		return
	}
}

func TestGetCommand(t *testing.T) {
	c := TestingCore("get_command")

	c.FetchCommands()

	c.AddCommand("test", "ls")

	// "Refetch" the commands config
	c.FetchCommands()

	exists, command := c.CommandExists("test")

	if !exists {
		t.Error("Expected command to exist, got", exists)
		return
	}

	if command.Command != "ls" {
		t.Error("Expected command command to be 'ls', got", command)
		return
	}
}

func TestGetCommandNotFound(t *testing.T) {
	c := TestingCore("get_command_not_found")

	c.FetchCommands()

	exists, _ := c.CommandExists("test")

	if exists {
		t.Error("Expected command to not exist, got", exists)
		return
	}
}

func TestEditCommand(t *testing.T) {
	c := TestingCore("edit_command")

	c.FetchCommands()

	c.AddCommand("test", "ls")

	// "Refetch" the commands config
	c.FetchCommands()

	c.UpdateCommand("test", "pwd")

	// "Refetch" the commands config
	c.FetchCommands()

	exists, command := c.CommandExists("test")

	if !exists {
		t.Error("Expected command to exist, got", exists)
		return
	}

	if command.Command != "pwd" {
		t.Error("Expected command command to be 'pwd', got", command)
		return
	}
}

func TestRenameCommand(t *testing.T) {
	c := TestingCore("rename_command")

	c.FetchCommands()

	c.AddCommand("test", "ls")

	// "Refetch" the commands config
	c.FetchCommands()

	c.RenameCommand("test", "test2")

	// "Refetch" the commands config
	c.FetchCommands()

	exists, command := c.CommandExists("test2")

	if !exists {
		t.Error("Expected command to exist, got", exists)
		return
	}

	if command.Command != "ls" {
		t.Error("Expected command command to be 'ls', got", command)
		return
	}
}
