package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/iskandervdh/spinup/config"
)

type Commands map[string]string

func (s *Core) getCommandsFilePath() string {
	return path.Join(s.config.GetConfigDir(), config.CommandsFileName)
}

func (s *Core) getCommands() (Commands, error) {
	commandsFileContent, err := os.ReadFile(s.getCommandsFilePath())

	if err != nil {
		s.cli.ErrorPrint("Error reading commands.json file:", err)
		return nil, err
	}

	var commands Commands
	err = json.Unmarshal(commandsFileContent, &commands)

	if err != nil {
		s.cli.ErrorPrint("Error parsing commands.json file:", err)
		return nil, err
	}

	return commands, nil
}

func (s *Core) getCommand(name string) (string, error) {
	command, exists := s.commands[name]

	if !exists {
		return "", fmt.Errorf("command '%s' not found", name)
	}

	return command, nil
}

func (s *Core) addCommand(name string, command string) {
	// Check if already exists
	for commandName := range s.commands {
		if commandName == name {
			s.cli.ErrorPrintf("Command '%s' already exists", name)
			return
		}
	}

	s.commands[name] = command

	updatedCommands, err := json.MarshalIndent(s.commands, "", "  ")

	if err != nil {
		s.cli.ErrorPrint("Error encoding projects to json:", err)
		return
	}

	err = os.WriteFile(s.getCommandsFilePath(), updatedCommands, 0644)

	if err != nil {
		s.cli.ErrorPrint("Error writing commands to file:", err)
		return
	}

	if !s.config.IsTesting() {
		s.cli.InfoPrintf("Added command '%s': %s", name, command)
	}
}

func (s *Core) addCommandInteractive() {
	name := s.cli.Input("Enter command name:")
	command := s.cli.Input("Enter command:")

	s.addCommand(name, command)
}

func (s *Core) removeCommand(name string) {
	if s.commands == nil {
		return
	}

	delete(s.commands, name)

	updatedCommands, err := json.MarshalIndent(s.commands, "", "  ")

	if err != nil {
		s.cli.ErrorPrint("Error encoding commands to json:", err)
		return
	}

	err = os.WriteFile(s.getCommandsFilePath(), updatedCommands, 0644)

	if err != nil {
		s.cli.ErrorPrint("Error writing commands to file:", err)
		return
	}

	if !s.config.IsTesting() {
		s.cli.InfoPrintf("Removed command '%s'", name)
	}
}

func (s *Core) removeCommandInteractive() {
	name, err, exited := s.cli.Selection("Select command to remove", s.getCommandNames())

	if err != nil {
		s.cli.ErrorPrint("Error selecting command:", err)
		return
	}

	if exited {
		return
	}

	if name == "" {
		s.cli.ErrorPrint("No command selected")
		return
	}

	if !s.cli.Confirm("Are you sure you want to remove command " + name + "?") {
		return
	}

	s.removeCommand(name)
}

func (s *Core) listCommands() {
	if s.commands == nil {
		return
	}

	fmt.Printf("%-20s %-30s\n", "Name", "Command")

	for commandName, command := range s.commands {
		fmt.Printf("%-20s %-30s\n", commandName, command)
	}
}

func (s *Core) handleCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: spinup command <add|remove|list> [args...]")
		return
	}

	commandName := os.Args[2]

	switch commandName {
	case "list", "ls":
		s.listCommands()
	case "add":
		if len(os.Args) == 3 {
			s.addCommandInteractive()
			return
		}

		if len(os.Args) < 5 {
			fmt.Printf("Usage: %s command add <name> <command>\n", config.ProgramName)
			return
		}

		s.addCommand(os.Args[3], os.Args[4])
	case "remove", "rm":
		if len(os.Args) == 3 {
			s.removeCommandInteractive()
			return
		}

		if len(os.Args) < 4 {
			fmt.Printf("Usage: %s command remove <name>\n", config.ProgramName)
			return
		}

		s.removeCommand(os.Args[3])
	default:
		fmt.Printf("Unknown subcommand '%s'\n", commandName)
		fmt.Println("Expected 'add' or 'list'")
	}
}