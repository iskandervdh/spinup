package core

import (
	"encoding/json"
	"os"

	"github.com/iskandervdh/spinup/common"
)

type Variables map[string]string

func (c *Core) AddVariable(name string, key string, value string) common.Msg {
	if c.projects == nil {
		return common.NewErrMsg("no projects found")
	}

	exists, project := c.ProjectExists(name)

	if !exists {
		return common.NewErrMsg("project '%s' does not exist", name)
	}

	// Check if the variable is already defined
	for variableKey := range project.Variables {
		if variableKey == key {
			return common.NewErrMsg("variable with name '%s' already exists", name)
		}
	}

	project.Variables[key] = value

	updatedProjects := c.projects
	updatedProjects[name] = project

	updatedVariables, err := json.MarshalIndent(updatedProjects, "", "  ")

	if err != nil {
		return common.NewErrMsg("error encoding projects to json: %s", err)
	}

	err = os.WriteFile(c.getProjectsFilePath(), updatedVariables, 0644)

	if err != nil {
		return common.NewErrMsg("error writing projects to file: %s", err)
	}

	return common.NewSuccessMsg("Added variable '%s' to project '%s' with value '%s'\n", key, name, value)

}

func (c *Core) RemoveVariable(name string, key string) common.Msg {
	if c.projects == nil {
		return common.NewErrMsg("no projects found")
	}

	exists, project := c.ProjectExists(name)

	if !exists {
		return common.NewErrMsg("project '%s' does not exist, nothing to remove", name)
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
	updatedProjects[name] = project

	updatedProjectConfig, err := json.MarshalIndent(updatedProjects, "", "  ")

	if err != nil {
		return common.NewErrMsg("Error encoding projects to json: %s", err)
	}

	err = os.WriteFile(c.getProjectsFilePath(), updatedProjectConfig, 0644)

	if err != nil {
		return common.NewErrMsg("Error writing projects to file: %s", err)
	}

	return common.NewSuccessMsg("Removed variable '%s' from project '%s'\n", key, name)
}
