package spinup

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/iskandervdh/spinup/config"
)

type Variables map[string]string

func (s *Spinup) addVariable(name string, key string, value string) error {
	if s.projects == nil {
		return fmt.Errorf("no projects found")
	}

	exists, project := s.projectExists(name)

	if !exists {
		return fmt.Errorf("project '%s' does not exist", name)
	}

	// Check if the variable is already defined
	for variableKey := range project.Variables {
		if variableKey == key {
			return fmt.Errorf("variable with name '%s' already exists", name)
		}
	}

	project.Variables[key] = value

	updatedProjects := s.projects
	updatedProjects[name] = project

	updatedVariables, err := json.MarshalIndent(updatedProjects, "", "  ")

	if err != nil {
		return fmt.Errorf("error encoding projects to json: %s", err)
	}

	err = os.WriteFile(s.getProjectsFilePath(), updatedVariables, 0644)

	if err != nil {
		return fmt.Errorf("error writing projects to file: %s", err)
	}

	if !s.config.IsTesting() {
		s.cli.InfoPrintf("Added variable '%s' to project '%s' with value '%s'\n", key, name, value)
	}

	return nil
}

func (s *Spinup) removeVariable(name string, key string) error {
	if s.projects == nil {
		return fmt.Errorf("no projects found")
	}

	exists, project := s.projectExists(name)

	if !exists {
		return fmt.Errorf("project '%s' does not exist, nothing to remove", name)
	}

	if project.Variables[key] == "" {
		return fmt.Errorf("variable '%s' does not exist", key)
	}

	variables := make(map[string]string)

	for variableKey, variableValue := range project.Variables {
		if variableKey == key {
			continue
		}

		variables[variableKey] = variableValue
	}

	project.Variables = variables
	updatedProjects := s.projects
	updatedProjects[name] = project

	updatedProjectConfig, err := json.MarshalIndent(updatedProjects, "", "  ")

	if err != nil {
		return fmt.Errorf("error encoding projects to json: %s", err)
	}

	err = os.WriteFile(s.getProjectsFilePath(), updatedProjectConfig, 0644)

	if err != nil {
		return fmt.Errorf("error writing projects to file: %s", err)
	}

	if !s.config.IsTesting() {
		s.cli.InfoPrintf("Removed variable '%s' from project '%s'\n", key, name)
	}

	return nil
}

func (s *Spinup) listVariables(name string) error {
	if s.projects == nil {
		return fmt.Errorf("no projects found")
	}

	exists, project := s.projectExists(name)

	if !exists {
		return fmt.Errorf("project '%s' does not exist", name)
	}

	fmt.Printf("%-20s %-30s\n", "Key", "Value")

	for key, value := range project.Variables {
		fmt.Printf("%-20s %-30s\n", key, value)
	}

	return nil
}

func (s *Spinup) handleVariable() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s variable <add|remove|list> [args...]\n", config.ProgramName)
		return
	}

	switch os.Args[2] {
	case "list", "ls":
		if len(os.Args) < 4 {
			fmt.Printf("Usage: %s variable list|ls <project>\n", config.ProgramName)
			return
		}

		err := s.listVariables(os.Args[3])

		if err != nil {
			s.cli.ErrorPrint("Error listing variables:", err)
		}
	case "add":
		if len(os.Args) < 6 {
			fmt.Printf("Usage: %s variable add <project> <key> <value>\n", config.ProgramName)
			return
		}

		s.addVariable(os.Args[3], os.Args[4], os.Args[5])
	case "remove", "rm":
		if len(os.Args) < 5 {
			fmt.Printf("Usage: %s variable remove|rm <project> <key>\n", config.ProgramName)
			return
		}

		err := s.removeVariable(os.Args[3], os.Args[4])

		if err != nil {
			s.cli.ErrorPrint("Error removing variable:", err)
		}
	}
}
