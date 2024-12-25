package core

import (
	"encoding/json"
	"os"

	"github.com/iskandervdh/spinup/common"
)

// variables is a map of variable names to their values.
type variables map[string]string

// Add a variable with the given key and value to the project with the given name.
func (c *Core) AddVariable(projectName string, key string, value string) common.Msg {
	if c.projects == nil {
		return common.NewErrMsg("no projects found")
	}

	exists, project := c.ProjectExists(projectName)

	if !exists {
		return common.NewErrMsg("project '%s' does not exist", projectName)
	}

	// Check if the variable is already defined
	for variableKey := range project.Variables {
		if variableKey == key {
			return common.NewErrMsg("variable with name '%s' already exists", projectName)
		}
	}

	project.Variables[key] = value

	updatedProjects := c.projects
	updatedProjects[projectName] = project

	updatedVariables, err := json.MarshalIndent(updatedProjects, "", "  ")

	if err != nil {
		return common.NewErrMsg("error encoding projects to json: %s", err)
	}

	err = os.WriteFile(c.getProjectsFilePath(), updatedVariables, 0644)

	if err != nil {
		return common.NewErrMsg("error writing projects to file: %s", err)
	}

	return common.NewSuccessMsg("Added variable '%s' to project '%s' with value '%s'\n", key, projectName, value)

}

// Remove the variable with the given key from the project with the given name.
func (c *Core) RemoveVariable(projectName string, key string) common.Msg {
	if c.projects == nil {
		return common.NewErrMsg("no projects found")
	}

	exists, project := c.ProjectExists(projectName)

	if !exists {
		return common.NewErrMsg("project '%s' does not exist, nothing to remove", projectName)
	}

	if project.Variables[key] == "" {
		return common.NewErrMsg("variable '%s' does not exist", key)
	}

	variables := make(map[string]string)

	for variableKey, variableValue := range project.Variables {
		if variableKey == key {
			continue
		}

		variables[variableKey] = variableValue
	}

	project.Variables = variables
	updatedProjects := c.projects
	updatedProjects[projectName] = project

	updatedProjectConfig, err := json.MarshalIndent(updatedProjects, "", "  ")

	if err != nil {
		return common.NewErrMsg("Error encoding projects to json: %s", err)
	}

	err = os.WriteFile(c.getProjectsFilePath(), updatedProjectConfig, 0644)

	if err != nil {
		return common.NewErrMsg("Error writing projects to file: %s", err)
	}

	return common.NewSuccessMsg("Removed variable '%s' from project '%s'\n", key, projectName)
}
