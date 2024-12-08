package spinup

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/iskandervdh/spinup/cli"
)

func (s *Spinup) createConfigDir() error {
	// Create config directory if it doesn't exist
	configDir := path.Dir(s.config.ConfigDir)
	err := os.MkdirAll(configDir, 0755)

	if err != nil {
		return err
	}

	return nil
}

func (s *Spinup) createProjectsConfigFile() error {
	projectFilePath := s.getProjectsFilePath()

	if _, err := os.Stat(projectFilePath); err == nil {
		cli.WarningPrintf(
			"%s file already exists at %s\nSkipping initialization...\n",
			s.config.ProjectsFileName,
			projectFilePath,
		)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error checking if %s file exists: %w", s.config.ProjectsFileName, err)
	}

	emptyProjects := Projects{}
	emptyData, err := json.MarshalIndent(emptyProjects, "", "  ")

	if err != nil {
		return fmt.Errorf("error encoding empty projects to json: %w", err)
	}

	err = os.WriteFile(projectFilePath, emptyData, 0644)

	if err != nil {
		return fmt.Errorf("error writing empty projects to file: %w", err)
	}

	cli.InfoPrint("Initialized empty projects.json file at ", projectFilePath)

	return nil
}

func (s *Spinup) createCommandsConfigFile() error {
	commandsFilePath := s.getCommandsFilePath()

	if _, err := os.Stat(commandsFilePath); err == nil {
		cli.WarningPrintf(
			"%s file already exists at %s\nSkipping initialization...\n",
			s.config.CommandsFileName,
			commandsFilePath,
		)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("error checking if commands.json file exists: %w", err)
	}

	emptyCommands := Commands{}
	emptyData, err := json.MarshalIndent(emptyCommands, "", "  ")

	if err != nil {
		return fmt.Errorf("error encoding empty commands to json: %w", err)
	}

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
	err = s.config.InitHosts()

	if err != nil {
		cli.ErrorPrint("Error initializing hosts file:", err)
		os.Exit(1)
	}
}
