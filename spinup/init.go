package spinup

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/iskandervdh/spinup/config"
)

func (s *Spinup) createConfigDir() error {
	// Create config directory if it doesn't exist
	configDir := path.Dir(s.configDirPath)
	err := os.MkdirAll(configDir, 0755)

	if err != nil {
		return err
	}

	return nil
}

func (s *Spinup) createProjectsConfigFile() error {
	emptyProjects := Projects{}
	emptyData, err := json.MarshalIndent(emptyProjects, "", "  ")

	if err != nil {
		return fmt.Errorf("error encoding empty projects to json: %w", err)
	}

	projectFilePath := s.getProjectsFilePath()
	err = os.WriteFile(projectFilePath, emptyData, 0644)

	if err != nil {
		return fmt.Errorf("error writing empty projects to file: %w", err)
	}

	cli.InfoPrint("Initialized empty projects.json file at ", projectFilePath)

	return nil
}

func (s *Spinup) createCommandsConfigFile() error {
	emptyCommands := Commands{}
	emptyData, err := json.MarshalIndent(emptyCommands, "", "  ")

	if err != nil {
		return fmt.Errorf("error encoding empty commands to json: %w", err)
	}

	commandsFilePath := s.getCommandsFilePath()
	err = os.WriteFile(commandsFilePath, emptyData, 0644)

	if err != nil {
		return fmt.Errorf("error writing empty commands to file: %w", err)
	}

	cli.InfoPrint("Initialized empty commands.json file at ", commandsFilePath)

	return nil
}

func (s *Spinup) init() {
	err := s.createConfigDir()

	if err != nil {
		cli.ErrorPrint("Error creating config directory: ", err)
		os.Exit(1)
	}

	err = s.createProjectsConfigFile()

	if err != nil {
		cli.ErrorPrint("Error creating projects.json file:", err)
		os.Exit(1)
	}

	err = s.createCommandsConfigFile()

	if err != nil {
		cli.ErrorPrint("Error creating commands.json file:", err)
		os.Exit(1)
	}

	s.requireSudo()
	err = config.InitHosts()

	if err != nil {
		cli.ErrorPrint("Error initializing hosts file:", err)
		os.Exit(1)
	}
}
