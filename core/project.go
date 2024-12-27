package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
)

// Project is a struct that represents a Project as stored in the projects.json file.
type Project struct {
	Domain        string        `json:"domain"`
	Port          int           `json:"port"`
	Commands      []string      `json:"commands"`
	Dir           *string       `json:"dir"`
	Variables     Variables     `json:"variables"`
	DomainAliases DomainAliases `json:"domainAliases"`
}

// Projects is a map of project names to their Projects.
type Projects map[string]Project

// Get the path to the projects.json file.
func (c *Core) getProjectsFilePath() string {
	return path.Join(c.config.GetConfigDir(), config.ProjectsFileName)
}

// Get the projects from the projects.json file.
func (c *Core) GetProjects() (Projects, error) {
	projectsFileContent, err := os.ReadFile(c.getProjectsFilePath())

	if err != nil {
		return nil, fmt.Errorf("error reading projects.json file: %s", err)
	}

	var projects Projects
	err = json.Unmarshal(projectsFileContent, &projects)

	if err != nil {
		return nil, fmt.Errorf("error parsing projects.json file: %s", err)
	}

	return projects, nil
}

// Check if a project with the given name exists. Returns the project if it exists.
func (c *Core) ProjectExists(name string) (bool, Project) {
	if c.projects == nil {
		return false, Project{}
	}

	project, exists := c.projects[name]

	return exists, project
}

// Add a project with the given name, domain, port and command names.
func (c *Core) AddProject(name string, domain string, port int, commandNames []string) common.Msg {
	c.RequireSudo()

	// Check if commands exist
	for _, commandName := range commandNames {
		_, exists := c.commands[commandName]

		if !exists {
			return common.NewErrMsg("Command %s does not exist", commandName)
		}
	}

	// Check if project already exists or if domain or port is already in use
	for projectName, project := range c.projects {
		if projectName == name {
			return common.NewErrMsg("Project '" + name + "' already exists")
		}

		if project.Domain == domain {
			return common.NewErrMsg("Project with domain '" + domain + "' already exists: " + projectName)

		}

		if project.Port == port {
			return common.NewErrMsg("Project with port " + strconv.Itoa(port) + " already exists: " + projectName)
		}
	}

	err := c.config.AddNginxConfig(name, domain, port)

	if err != nil {
		return common.NewErrMsg(fmt.Sprintln("Error trying to create nginx config file", err))
	}

	err = c.config.AddDomain(domain)

	if err != nil {
		// Remove nginx config file if adding domain to hosts file fails
		c.config.RemoveNginxConfig(name)

		return common.NewErrMsg(fmt.Sprintln("Error trying to add domain to hosts file", err))
	}

	newProject := Project{
		Domain:        domain,
		Port:          port,
		Commands:      commandNames,
		Variables:     Variables{},
		DomainAliases: DomainAliases{},
	}

	c.projects[name] = newProject

	updatedProjectsConfig, err := json.MarshalIndent(c.projects, "", "  ")

	if err != nil {
		return common.NewErrMsg(fmt.Sprintln("Error encoding projects to json:", err))
	}

	err = os.WriteFile(c.getProjectsFilePath(), updatedProjectsConfig, 0644)

	if err != nil {
		return common.NewErrMsg(fmt.Sprintln("Error writing projects to file:", err))
	}

	return common.NewSuccessMsg(fmt.Sprintf("Added project '%s' with domain '%s', port %d and commands %s", name, domain, port, commandNames))
}

// Remove the project with the given name.
func (c *Core) RemoveProject(name string) common.Msg {
	c.RequireSudo()

	exists, _ := c.ProjectExists(name)

	if !exists {
		return common.NewErrMsg("Project '" + name + "' does not exist, nothing to remove")
	}

	err := c.config.RemoveNginxConfig(name)

	if err != nil {
		return common.NewErrMsg("Could not remove nginx config file: " + err.Error())
	}

	err = c.config.RemoveDomain(c.projects[name].Domain)

	if err != nil {
		// Remove nginx config file if adding domain to hosts file fails
		c.config.RemoveNginxConfig(name)

		return common.NewErrMsg("Error trying to remove domain from hosts file: " + err.Error())
	}

	var updatedProjects Projects = make(map[string]Project)

	for projectName, project := range c.projects {
		if projectName == name {
			continue
		}

		updatedProjects[projectName] = project
	}

	updatedProjectsConfig, err := json.MarshalIndent(updatedProjects, "", "  ")

	if err != nil {
		return common.NewErrMsg("Error encoding projects to json: " + err.Error())
	}

	err = os.WriteFile(c.getProjectsFilePath(), updatedProjectsConfig, 0644)

	if err != nil {
		return common.NewErrMsg("Error writing projects to file: " + err.Error())
	}

	return common.NewSuccessMsg(fmt.Sprintf("Removed project '%s'", name))
}

// Update the project with the given name to the given domain, port and command names.
func (c *Core) UpdateProject(name string, domain string, port int, commandNames []string) common.Msg {
	c.RequireSudo()

	exists, project := c.ProjectExists(name)

	if !exists {
		return common.NewErrMsg("Project '%s' does not exist", name)
	}

	// Check if commands exist
	for _, commandName := range commandNames {
		_, exists := c.commands[commandName]

		if !exists {
			return common.NewErrMsg("Command '%s' does not exist", commandName)
		}
	}

	// Check if domain or port is already in use
	for projectName, p := range c.projects {
		if projectName == name {
			continue
		}

		if p.Domain == domain {
			return common.NewErrMsg("Project with domain '%s' already exists: %s", domain, projectName)
		}

		if p.Port == port {
			return common.NewErrMsg("Project with port %d already exists: %s", port, projectName)
		}
	}

	err := c.config.UpdateNginxConfig(name, domain, port)

	if err != nil {
		return common.NewErrMsg("Error trying to update nginx config file: %s", err)
	}

	err = c.config.UpdateHost(project.Domain, domain)

	if err != nil {
		// Undo editing of nginx config file if adding new domain to hosts file fails
		c.config.UpdateNginxConfig(name, project.Domain, project.Port)

		return common.NewErrMsg("Error trying to update domain in hosts file: %s", err)
	}

	project.Domain = domain
	project.Port = port
	project.Commands = commandNames

	c.projects[name] = project

	updatedProjectsConfig, err := json.MarshalIndent(c.projects, "", "  ")

	if err != nil {
		return common.NewErrMsg("Error encoding projects to json: %s", err)
	}

	err = os.WriteFile(c.getProjectsFilePath(), updatedProjectsConfig, 0644)

	if err != nil {
		return common.NewErrMsg("Error writing projects to file: %s", err)
	}

	return common.NewSuccessMsg("Updated project '%s' with domain '%s', port %d and commands %s", name, domain, port, commandNames)
}

// Rename the project with the given old name to the given new name.
func (c *Core) RenameProject(oldName string, newName string) common.Msg {
	exists, project := c.ProjectExists(oldName)

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

	c.projects[newName] = project

	delete(c.projects, oldName)

	updatedProjectsConfig, err := json.MarshalIndent(c.projects, "", "  ")

	if err != nil {
		return common.NewErrMsg("Error encoding projects to json: %s", err)
	}

	err = os.WriteFile(c.getProjectsFilePath(), updatedProjectsConfig, 0644)

	if err != nil {
		return common.NewErrMsg("Error writing projects to file: %s", err)
	}

	return common.NewSuccessMsg("Renamed project '%s' to '%s'", oldName, newName)
}

// Add a command to the project with the given name.
func (c *Core) AddCommandToProject(projectName string, commandName string) common.Msg {
	exists, project := c.ProjectExists(projectName)

	if !exists {
		return common.NewErrMsg("Project '%s' does not exist", projectName)
	}

	_, exists = c.commands[commandName]

	if !exists {
		return common.NewErrMsg("Command '%s' does not exist", commandName)
	}

	for _, command := range project.Commands {
		if command == commandName {
			return common.NewErrMsg("Command '%s' already exists in project '%s'", commandName, projectName)
		}
	}

	project.Commands = append(project.Commands, commandName)

	c.projects[projectName] = project

	updatedProjectsConfig, err := json.MarshalIndent(c.projects, "", "  ")

	if err != nil {
		return common.NewErrMsg("Error encoding projects to json:", err)
	}

	err = os.WriteFile(c.getProjectsFilePath(), updatedProjectsConfig, 0644)

	if err != nil {
		return common.NewErrMsg("Error writing projects to file:", err)
	}

	return common.NewSuccessMsg("Added command '%s' to project '%s'", commandName, projectName)
}

// Remove the command with the given commandName from the project with the given projectName.
func (c *Core) RemoveCommandFromProject(projectName string, commandName string) common.Msg {
	exists, project := c.ProjectExists(projectName)

	if !exists {
		return common.NewErrMsg("Project '%s' does not exist", projectName)
	}

	for i, command := range project.Commands {
		if command == commandName {
			project.Commands = append(project.Commands[:i], project.Commands[i+1:]...)

			c.projects[projectName] = project

			updatedProjectsConfig, err := json.MarshalIndent(c.projects, "", "  ")

			if err != nil {
				return common.NewErrMsg("Error encoding projects to json: %s", err)
			}

			err = os.WriteFile(c.getProjectsFilePath(), updatedProjectsConfig, 0644)

			if err != nil {
				return common.NewErrMsg("Error writing projects to file: %s", err)
			}

			return common.NewSuccessMsg("Removed command '%s' from project '%s'", commandName, projectName)
		}
	}

	return common.NewInfoMsg("Command '%s' not found in project '%s'. Nothing to remove.", commandName, projectName)
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

		project.Dir = &cwd
	} else {
		// Check if dir exists as a directory
		info, err := os.Stat(*dir)

		if err != nil {
			return common.NewErrMsg("Directory '%s' does not exist: %s", *dir, err)
		}

		if !info.IsDir() {
			return common.NewErrMsg("'%s' is not a directory", *dir)
		}

		project.Dir = dir
	}

	c.projects[projectName] = project

	updatedProjectsConfig, err := json.MarshalIndent(c.projects, "", "  ")

	if err != nil {
		return common.NewErrMsg("Error encoding projects to json: %s", err)
	}

	err = os.WriteFile(c.getProjectsFilePath(), updatedProjectsConfig, 0644)

	if err != nil {
		return common.NewErrMsg("Error writing projects to file: %s", err)
	}

	return common.NewSuccessMsg("Set directory to '%s' for project '%s'", *project.Dir, projectName)
}

// Get the directory for the project with the given name.
func (c *Core) GetProjectDir(projectName string) common.Msg {
	exists, project := c.ProjectExists(projectName)

	if !exists {
		return common.NewErrMsg("Project '%s' does not exist", projectName)
	}

	if project.Dir == nil {
		return common.NewErrMsg("Project '%s' does not have a directory set", projectName)
	}

	return common.NewRegularMsg(*project.Dir)
}
