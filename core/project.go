package core

import (
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/database/sqlc"
)

// Project is a struct that represents a project and its linked structs.
type Project struct {
	sqlc.Project
	Commands      []Command
	Variables     []Variable
	DomainAliases []DomainAlias
}

// Projects is a map of project names to their Projects.
type Projects []Project

func (c *Core) getProjectWithInfo(project sqlc.Project) (Project, error) {
	projectCommands, err := c.dbQueries.GetProjectCommands(c.dbContext, project.ID)

	if err != nil {
		return Project{}, fmt.Errorf("error getting project commands: %s", err)
	}

	projectVariables, err := c.dbQueries.GetProjectVariables(c.dbContext, project.ID)

	if err != nil {
		return Project{}, fmt.Errorf("error getting project variables: %s", err)
	}

	projectDomainAliases, err := c.dbQueries.GetProjectDomainAliases(c.dbContext, project.ID)

	if err != nil {
		return Project{}, fmt.Errorf("error getting project domain aliases: %s", err)
	}

	return Project{
		Project:       project,
		Commands:      projectCommands,
		Variables:     projectVariables,
		DomainAliases: projectDomainAliases,
	}, nil
}

func (c *Core) FetchProjects() error {
	projects, err := c.dbQueries.GetProjects(c.dbContext)

	if err != nil {
		return fmt.Errorf("error getting commands: %s", err)
	}

	var projectsWithInfo Projects

	for _, project := range projects {
		projectWithInfo, err := c.getProjectWithInfo(project)

		if err != nil {
			return err
		}

		projectsWithInfo = append(projectsWithInfo, projectWithInfo)
	}

	c.projects = projectsWithInfo

	return nil
}

// Get the projects from the database.
func (c *Core) GetProjects() (Projects, error) {
	if c.projects == nil {
		err := c.FetchProjects()

		if err != nil {
			return nil, err
		}
	}

	return c.projects, nil
}

// Check if a project with the given name exists. Returns the project if it exists.
func (c *Core) ProjectExists(name string) (bool, Project) {
	if c.projects == nil {
		return false, Project{}
	}

	index := slices.IndexFunc(c.projects, func(project Project) bool {
		return project.Name == name
	})

	if index == -1 {
		return false, Project{}
	}

	return true, c.projects[index]
}

// Add a project with the given name, port and command names.
func (c *Core) AddProject(name string, port int64, commandNames []string) common.Msg {
	// Check if commands exist
	commandIDs := make([]int64, 0, len(commandNames))

	for _, commandName := range commandNames {
		exists, command := c.CommandExists(commandName)

		if !exists {
			c.sendMsg(common.NewErrMsg("Command '" + commandName + "' does not exist"))
		}

		commandIDs = append(commandIDs, command.ID)
	}

	// Check if project already exists or the port is already in use
	for _, project := range c.projects {
		if project.Name == name {
			return common.NewErrMsg("Project '" + name + "' already exists")
		}

		if project.Port == port {
			return common.NewErrMsg("Project with port " + strconv.FormatInt(port, 10) + " already exists: " + project.Name)
		}
	}

	err := c.config.AddNginxConfig(name, port)

	if err != nil {
		return common.NewErrMsg(fmt.Sprintln("Error trying to create nginx config file", err))
	}

	project, err := c.dbQueries.CreateProject(c.dbContext, sqlc.CreateProjectParams{
		Name: name,
		Port: int64(port),
	})

	if err != nil {
		return common.NewErrMsg(fmt.Sprintf("Error adding project to database: %s", err))
	}

	for _, commandID := range commandIDs {
		err = c.dbQueries.CreateProjectCommand(c.dbContext, sqlc.CreateProjectCommandParams{
			ProjectID: project.ID,
			CommandID: commandID,
		})

		if err != nil {
			return common.NewErrMsg(fmt.Sprintf("Error adding commands to project in database: %s", err))
		}
	}

	return common.NewSuccessMsg(fmt.Sprintf("Added project '%s'", name))
}

// Remove the project with the given name.
func (c *Core) RemoveProject(name string) common.Msg {
	exists, _ := c.ProjectExists(name)

	if !exists {
		return common.NewErrMsg("Project '" + name + "' does not exist, nothing to remove")
	}

	err := c.config.RemoveNginxConfig(name)

	if err != nil {
		return common.NewErrMsg("Could not remove nginx config file: " + err.Error())
	}

	c.dbQueries.DeleteProject(c.dbContext, name)

	return common.NewSuccessMsg(fmt.Sprintf("Removed project '%s'", name))
}

// Update the project with the given name to the given port and command names.
func (c *Core) UpdateProject(name string, port int64, commandNames []string) common.Msg {
	exists, project := c.ProjectExists(name)

	if !exists {
		return common.NewErrMsg("Project '%s' does not exist", name)
	}

	// Check if commands exist
	for _, commandName := range commandNames {
		exists, _ := c.CommandExists(commandName)

		if !exists {
			return common.NewErrMsg("Command '%s' does not exist", commandName)
		}
	}

	// Check if port is already in use by another project
	for _, project := range c.projects {
		if project.Name == name {
			continue
		}

		if project.Port == port {
			return common.NewErrMsg("Project with port %d already exists: %s", port, project.Name)
		}
	}

	err := c.config.UpdateNginxConfig(name, port)

	if err != nil {
		return common.NewErrMsg("Error trying to update nginx config file: %s", err)
	}

	c.dbQueries.UpdateProject(c.dbContext, sqlc.UpdateProjectParams{
		Name: name,
		Port: port,
		Dir: sql.NullString{
			String: project.Dir.String,
			Valid:  project.Dir.Valid,
		},
	})

	return common.NewSuccessMsg("Updated project '%s' with domain '%s', port %d and commands %s", name, port, commandNames)
}

// Rename the project with the given old name to the given new name.
func (c *Core) RenameProject(oldName string, newName string) common.Msg {
	exists, _ := c.ProjectExists(oldName)

	if !exists {
		return common.NewErrMsg("Project '%s' does not exist", oldName)
	}

	newNameProjectExists, _ := c.ProjectExists(newName)

	if newNameProjectExists {
		return common.NewErrMsg("Project '%s' already exists", newName)
	}

	err := c.config.RenameNginxConfig(oldName, newName)

	if err != nil {
		return common.NewErrMsg("Error trying to rename nginx config file: %s", err)
	}

	err = c.dbQueries.RenameProject(c.dbContext, sqlc.RenameProjectParams{
		Name:   newName,
		Name_2: oldName,
	})

	if err != nil {
		return common.NewErrMsg("Error renaming project in database: %s", err)
	}

	return common.NewSuccessMsg("Renamed project '%s' to '%s'", oldName, newName)
}

// Add a command to the project with the given name.
func (c *Core) AddCommandToProject(projectName string, commandName string) common.Msg {
	exists, project := c.ProjectExists(projectName)

	if !exists {
		return common.NewErrMsg("Project '%s' does not exist", projectName)
	}

	exists, command := c.CommandExists(commandName)

	if !exists {
		return common.NewErrMsg("Command '%s' does not exist", commandName)
	}

	for _, projectCommand := range project.Commands {
		if projectCommand.ID == command.ID {
			return common.NewErrMsg("Command '%s' already exists in project '%s'", commandName, projectName)
		}
	}

	c.dbQueries.CreateProjectCommand(c.dbContext, sqlc.CreateProjectCommandParams{
		ProjectID: project.ID,
		CommandID: command.ID,
	})

	return common.NewSuccessMsg("Added command '%s' to project '%s'", commandName, projectName)
}

// Remove the command with the given commandName from the project with the given projectName.
func (c *Core) RemoveCommandFromProject(projectName string, commandName string) common.Msg {
	exists, project := c.ProjectExists(projectName)

	if !exists {
		return common.NewErrMsg("Project '%s' does not exist", projectName)
	}

	exists, command := c.CommandExists(commandName)

	if !exists {
		return common.NewErrMsg("Command '%s' does not exist", commandName)
	}

	err := c.dbQueries.DeleteCommandsProjects(c.dbContext, sqlc.DeleteCommandsProjectsParams{
		CommandID: command.ID,
		ProjectID: project.ID,
	})

	if err != nil {
		return common.NewErrMsg("Error removing command from project: %s", err)
	}

	return common.NewInfoMsg("Removed command '%s' from project '%s'", commandName, projectName)
}

// Set the directory for the project with the given name.
func (c *Core) SetProjectDir(projectName string, dir *string) common.Msg {
	exists, project := c.ProjectExists(projectName)

	if !exists {
		return common.NewErrMsg("Project '%s' does not exist", projectName)
	}

	if dir == nil {
		cwd, err := os.Getwd()

		if err != nil {
			return common.NewErrMsg("Error getting current working directory: %s", err)
		}

		project.Dir = sql.NullString{
			String: cwd,
			Valid:  true,
		}
	} else {
		// Check if dir exists as a directory
		info, err := os.Stat(*dir)

		if err != nil {
			return common.NewErrMsg("Directory '%s' does not exist: %s", *dir, err)
		}

		if !info.IsDir() {
			return common.NewErrMsg("'%s' is not a directory", *dir)
		}

		project.Dir = sql.NullString{
			String: *dir,
			Valid:  true,
		}
	}

	c.dbQueries.SetProjectDir(c.dbContext, sqlc.SetProjectDirParams{
		Dir: project.Dir,
		ID:  project.ID,
	})

	return common.NewSuccessMsg("Set directory to '%s' for project '%s'", project.Dir.String, projectName)
}

// Get the directory for the project with the given name.
func (c *Core) GetProjectDir(projectName string) common.Msg {
	exists, project := c.ProjectExists(projectName)

	if !exists {
		return common.NewErrMsg("Project '%s' does not exist", projectName)
	}

	if !project.Dir.Valid {
		return common.NewErrMsg("Project '%s' does not have a directory set", projectName)
	}

	return common.NewRegularMsg(project.Dir.String)
}
