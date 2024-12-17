package core

import (
	"sort"
	"testing"
)

func TestGetCommandNames(t *testing.T) {
	s := TestingCore("get_command_names")

	commandNames := s.GetCommandNames()

	if len(commandNames) != 0 {
		t.Error("Expected no command names, got", len(commandNames))
		return
	}
}

func TestAddCommand(t *testing.T) {
	s := TestingCore("add_command")

	s.GetCommandsConfig()

	s.AddCommand("test", "ls")

	// "Refetch" the commands config
	s.GetCommandsConfig()

	commandNames := s.GetCommandNames()

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
	s := TestingCore("add_command_duplicate")

	s.GetCommandsConfig()

	s.AddCommand("test", "ls")

	// "Refetch" the commands config
	s.GetCommandsConfig()

	s.AddCommand("test", "ls")

	// "Refetch" the commands config
	s.GetCommandsConfig()

	commandNames := s.GetCommandNames()

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
	s := TestingCore("add_multiple_commands")

	s.GetCommandsConfig()

	s.AddCommand("test", "ls")
	s.AddCommand("test2", "ls")

	// "Refetch" the commands config
	s.GetCommandsConfig()

	commandNames := s.GetCommandNames()
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
	s := TestingCore("remove_command")

	s.GetCommandsConfig()

	s.AddCommand("test", "ls")

	// "Refetch" the commands config
	s.GetCommandsConfig()

	s.RemoveCommand("test")

	// "Refetch" the commands config
	s.GetCommandsConfig()

	commandNames := s.GetCommandNames()

	if len(commandNames) != 0 {
		t.Error("Expected no command names, got", len(commandNames))
		return
	}
}

func TestGetCommand(t *testing.T) {
	s := TestingCore("get_command")

	s.GetCommandsConfig()

	s.AddCommand("test", "ls")

	// "Refetch" the commands config
	s.GetCommandsConfig()

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
	s := TestingCore("get_command_not_found")

	s.GetCommandsConfig()

	_, err := s.getCommand("test")

	if err == nil {
		t.Error("Expected error, got nil")
		return
	}
}

// func TestListCommands(t *testing.T) {
// 	s := TestingCore("list_commands")

// 	s.GetCommandsConfig()

// 	s.AddCommand("test", "ls")
// 	s.AddCommand("test2", "ls")

// 	// "Refetch" the commands config
// 	s.GetCommandsConfig()

// 	s.listCommands()
// }

// func TestListCommandsNoCommands(t *testing.T) {
// 	s := TestingCore("list_commands_no_commands")

// 	s.listCommands()
// }
