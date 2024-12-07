package spinup

import (
	"testing"
)

func TestGetCommandNames(t *testing.T) {
	s := TestingSpinup("get_command_names")

	commandNames := s.getCommandNames()

	if len(commandNames) != 0 {
		t.Error("Expected no command names, got", len(commandNames))
		return
	}
}

func TestAddCommand(t *testing.T) {
	s := TestingSpinup("add_command")

	s.getCommandsConfig()

	s.addCommand("test", "ls")

	// "Refetch" the commands config
	s.getCommandsConfig()

	commandNames := s.getCommandNames()

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
	s := TestingSpinup("add_command_duplicate")

	s.getCommandsConfig()

	s.addCommand("test", "ls")

	// "Refetch" the commands config
	s.getCommandsConfig()

	s.addCommand("test", "ls")

	// "Refetch" the commands config
	s.getCommandsConfig()

	commandNames := s.getCommandNames()

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
	s := TestingSpinup("add_multiple_commands")

	s.getCommandsConfig()

	s.addCommand("test", "ls")
	s.addCommand("test2", "ls")

	// "Refetch" the commands config
	s.getCommandsConfig()

	commandNames := s.getCommandNames()

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
	s := TestingSpinup("remove_command")

	s.getCommandsConfig()

	s.addCommand("test", "ls")

	// "Refetch" the commands config
	s.getCommandsConfig()

	s.removeCommand("test")

	// "Refetch" the commands config
	s.getCommandsConfig()

	commandNames := s.getCommandNames()

	if len(commandNames) != 0 {
		t.Error("Expected no command names, got", len(commandNames))
		return
	}
}

func TestGetCommand(t *testing.T) {
	s := TestingSpinup("get_command")

	s.getCommandsConfig()

	s.addCommand("test", "ls")

	// "Refetch" the commands config
	s.getCommandsConfig()

	command, err := s.getCommand("test")

	if err != nil {
		t.Error("Expected command to be found, got nil")
		return
	}

	if command != "ls" {
		t.Error("Expected command command to be 'ls', got", command)
		return
	}
}

func TestGetCommandNotFound(t *testing.T) {
	s := TestingSpinup("get_command_not_found")

	s.getCommandsConfig()

	_, err := s.getCommand("test")

	if err == nil {
		t.Error("Expected error, got nil")
		return
	}
}
