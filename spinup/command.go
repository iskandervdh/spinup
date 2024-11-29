package spinup

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/iskandervdh/spinup/config"
)

type Commands map[string]string

func (s *Spinup) getCommandsFilePath() string {
	return path.Join(s.configDirPath, config.CommandsFileName)
}

func (s *Spinup) getCommands() (Commands, error) {
	commandsFileContent, err := os.ReadFile(s.getCommandsFilePath())

	if err != nil {
		fmt.Println("Error reading commands.json file:", err)
		return nil, err
	}

	var commands Commands
	err = json.Unmarshal(commandsFileContent, &commands)

	if err != nil {
		fmt.Println("Error parsing commands.json file:", err)
		return nil, err
	}

	return commands, nil
}

func (s *Spinup) getCommand(name string) (string, error) {
	command, exists := s.commands[name]

	if !exists {
		return "", fmt.Errorf("command '%s' not found", name)
	}

	return command, nil
}

func (s *Spinup) addCommand(name string, command string) {
	if len(os.Args) < 4 {
		fmt.Println("Usage: spinup command add <name> <command>")
		return
	}

	// Check if already exists
	for commandName := range s.commands {
		if commandName == name {
			fmt.Printf("Command '%s' already exists\n", name)
			return
		}
	}

	s.commands[name] = command

	updatedCommands, err := json.MarshalIndent(s.commands, "", "  ")

	if err != nil {
		fmt.Println("Error encoding projects to json:", err)
		return
	}

	os.WriteFile(s.getCommandsFilePath(), updatedCommands, 0644)

	fmt.Printf("Added command '%s': %s\n", name, command)
}

func (s *Spinup) removeCommand(name string) {
	if s.commands == nil {
		return
	}

	delete(s.commands, name)

	updatedCommands, err := json.MarshalIndent(s.commands, "", "  ")

	if err != nil {
		fmt.Println("Error encoding commands to json:", err)
		return
	}

	os.WriteFile(s.getCommandsFilePath(), updatedCommands, 0644)

	fmt.Printf("Removed command '%s'\n", name)
}

func (s *Spinup) listCommands() {
	if s.commands == nil {
		return
	}

	fmt.Printf("%-20s %-30s\n", "Name", "Command")

	for commandName, command := range s.commands {
		fmt.Printf("%-20s %-30s\n", commandName, command)
	}
}

func (s *Spinup) handleCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: spinup command <add|remove|list> [args...]")
		return
	}

	commandName := os.Args[2]

	switch commandName {
	case "list", "ls":
		s.listCommands()
	case "add":
		s.addCommand(os.Args[3], os.Args[4])
	case "remove", "rm":
		s.removeCommand(os.Args[3])
	default:
		fmt.Printf("Unknown subcommand '%s'\n", commandName)
		fmt.Println("Expected 'add' or 'list'")
	}
}
