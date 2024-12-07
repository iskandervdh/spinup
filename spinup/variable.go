package spinup

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/iskandervdh/spinup/config"
)

type Variables map[string]string

func (s *Spinup) addVariable(name string, key string, value string) {
	if s.projects == nil {
		return
	}

	exists, project := s.projectExists(name)

	if !exists {
		fmt.Printf("Project '%s' does not exist\n", name)
		return
	}

	// Check if the variable is already defined
	for variableKey := range project.Variables {
		if variableKey == key {
			fmt.Printf("Variable with name '%s' already exists.\n", name)
			return
		}
	}

	project.Variables[key] = value

	updatedProjects := s.projects
	updatedProjects[name] = project

	updatedVariables, err := json.MarshalIndent(updatedProjects, "", "  ")

	if err != nil {
		fmt.Println("Error encoding projects to json:", err)
		return
	}

	err = os.WriteFile(s.getProjectsFilePath(), updatedVariables, 0644)

	if err != nil {
		fmt.Println("Error writing projects to file:", err)
		return
	}

	fmt.Printf("Added variable '%s' to project '%s' with value '%s'\n", key, name, value)
}

func (s *Spinup) removeVariable(name string, key string) {
	if s.projects == nil {
		return
	}

	exists, project := s.projectExists(name)

	if !exists {
		fmt.Printf("Project '%s' does not exist, nothing to remove\n", name)
		return
	}

	if project.Variables[key] == "" {
		fmt.Printf("Variable '%s' does not exist\n", key)
		return
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
		fmt.Println("Error encoding projects to json:", err)
		return
	}

	err = os.WriteFile(s.getProjectsFilePath(), updatedProjectConfig, 0644)

	if err != nil {
		fmt.Println("Error writing projects to file:", err)
		return
	}

	fmt.Printf("Removed variable '%s' from project '%s'\n", key, name)
}

func (s *Spinup) listVariables(name string) {
	if s.projects == nil {
		return
	}

	exists, project := s.projectExists(name)

	if !exists {
		fmt.Printf("Project '%s' does not exist\n", name)
		return
	}

	fmt.Printf("%-20s %-30s\n", "Key", "Value")

	for key, value := range project.Variables {
		fmt.Printf("%-20s %-30s\n", key, value)
	}
}

func (s *Spinup) handleVariable() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s variable <add|remove|list> [args...]\n", config.ProgramName)
		return
	}

	switch os.Args[2] {
	case "list", "ls":
		s.listVariables(os.Args[3])
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

		s.removeVariable(os.Args[3], os.Args[4])
	}
}
