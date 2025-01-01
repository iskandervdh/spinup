package core

import (
	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/database/sqlc"
)

type Variable = sqlc.Variable

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
	for _, variable := range project.Variables {
		if variable.Name == key {
			return common.NewErrMsg("variable with name '%s' already exists", projectName)
		}
	}

	err := c.dbQueries.CreateVariable(c.dbContext, sqlc.CreateVariableParams{
		ProjectID: project.ID,
		Name:      key,
		Value:     value,
	})

	if err != nil {
		return common.NewErrMsg("Error creating variable: %s", err)
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

	err := c.dbQueries.DeleteVariable(c.dbContext, sqlc.DeleteVariableParams{
		ProjectID: project.ID,
		Name:      key,
	})

	if err != nil {
		return common.NewErrMsg("Error deleting variable: %s", err)
	}

	return common.NewSuccessMsg("Removed variable '%s' from project '%s'\n", key, projectName)
}
