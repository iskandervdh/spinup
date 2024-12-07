package spinup

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/iskandervdh/spinup/cli"
	"github.com/iskandervdh/spinup/config"
)

type Project struct {
	Domain    string    `json:"domain"`
	Port      int       `json:"port"`
	Commands  []string  `json:"commands"`
	Dir       *string   `json:"dir"`
	Variables Variables `json:"variables"`
}

type Projects map[string]Project

func (s *Spinup) getProjectsFilePath() string {
	return path.Join(s.configDirPath, config.ProjectsFileName)
}

func (s *Spinup) getProjects() (Projects, error) {
	projectsFileContent, err := os.ReadFile(s.getProjectsFilePath())

	if err != nil {
		fmt.Println("Error reading projects.json file:", err)
		return nil, err
	}

	var projects Projects
	err = json.Unmarshal(projectsFileContent, &projects)

	if err != nil {
		fmt.Println("Error parsing projects.json file:", err)
		return nil, err
	}

	return projects, nil
}

func (s *Spinup) projectExists(name string) (bool, Project) {
	if s.projects == nil {
		return false, Project{}
	}

	project, exists := s.projects[name]

	return exists, project
}

func (s *Spinup) _addProject(name string, domain string, port int, commandNames []string) tea.Msg {
	// Check if commands exist
	for _, commandName := range commandNames {
		_, exists := s.commands[commandName]

		if !exists {
			fmt.Println("Command", commandName, "does not exist")
			return cli.ErrMsg("Command " + commandName + " does not exist")
		}
	}

	// Check if project already exists or if domain or port is already in use
	for projectName, project := range s.projects {
		if projectName == name {
			return cli.ErrMsg("Project '" + name + "' already exists")
		}

		if project.Domain == domain {
			return cli.ErrMsg("Project with domain '" + domain + "' already exists: " + projectName)

		}

		if project.Port == port {
			return cli.ErrMsg("Project with port " + strconv.Itoa(port) + " already exists: " + projectName)
		}
	}

	err := config.AddNginxConfig(name, domain, port)

	if err != nil {
		return cli.ErrMsg(fmt.Sprintln("Error trying to create nginx config file", err))
	}

	err = config.AddHost(domain)

	if err != nil {
		// Remove nginx config file if adding domain to hosts file fails
		config.RemoveNginxConfig(name)

		return cli.ErrMsg(fmt.Sprintln("Error trying to add domain to hosts file", err))
	}

	newProject := Project{
		Domain:    domain,
		Port:      port,
		Commands:  commandNames,
		Variables: make(map[string]string),
	}

	s.projects[name] = newProject

	updatedProjectsConfig, err := json.MarshalIndent(s.projects, "", "  ")

	if err != nil {
		return cli.ErrMsg(fmt.Sprintln("Error encoding projects to json:", err))
	}

	err = os.WriteFile(s.getProjectsFilePath(), updatedProjectsConfig, 0644)

	if err != nil {
		return cli.ErrMsg(fmt.Sprintln("Error writing projects to file:", err))
	}

	return cli.DoneMsg(fmt.Sprintf("Added project '%s' with domain '%s' and port %d", name, domain, port))
}

func (s *Spinup) addProject(name string, domain string, port int, commandNames []string) {
	s.requireSudo()

	if s.projects == nil {
		return
	}

	cli.Loading(fmt.Sprintf("Adding project %s...", name),
		func() tea.Msg {
			return s._addProject(name, domain, port, commandNames)
		},
	)
}

func (s *Spinup) addProjectInteractive() {
	name := cli.Input("Project name:")
	domain := cli.Input("Domain:")
	port := cli.Input("Port:")

	portInt, err := strconv.Atoi(port)

	if err != nil {
		fmt.Println("Port must be an integer")
		return
	}

	selectedCommands := cli.Question("Commands", s.getCommandNames())

	s.addProject(name, domain, portInt, selectedCommands)
}

func (s *Spinup) _removeProject(name string) tea.Msg {
	exists, _ := s.projectExists(name)

	if !exists {
		return cli.ErrMsg("Project '" + name + "' does not exist, nothing to remove")
	}

	err := config.RemoveNginxConfig(name)

	if err != nil {
		return cli.ErrMsg("Could not remove nginx config file: " + err.Error())
	}

	err = config.RemoveHost(s.projects[name].Domain)

	if err != nil {
		// Remove nginx config file if adding domain to hosts file fails
		config.RemoveNginxConfig(name)

		return cli.ErrMsg("Error trying to remove domain from hosts file: " + err.Error())

	}

	var updatedProjects Projects = make(map[string]Project)

	for projectName, project := range s.projects {
		if projectName == name {
			continue
		}

		updatedProjects[projectName] = project
	}

	updatedProjectsConfig, err := json.MarshalIndent(updatedProjects, "", "  ")

	if err != nil {
		return cli.ErrMsg("Error encoding projects to json: " + err.Error())
	}

	err = os.WriteFile(s.getProjectsFilePath(), updatedProjectsConfig, 0644)

	if err != nil {
		return cli.ErrMsg("Error writing projects to file: " + err.Error())
	}

	return cli.DoneMsg(fmt.Sprintf("Removed project '%s'", name))
}

func (s *Spinup) removeProject(name string) {
	s.requireSudo()

	if s.projects == nil {
		return
	}

	cli.Loading("Adding project...",
		func() tea.Msg {
			return s._removeProject(name)
		},
	)
}

func (s *Spinup) removeProjectInteractive() {
	name := cli.Selection("What project do you want to remove?", s.getProjectNames())

	if name == "" {
		fmt.Println("No project selected")
		return
	}

	if !cli.Confirm("Are you sure you want to remove project " + name + "?") {
		return
	}

	s.removeProject(name)
}

func (s *Spinup) listProjects() {
	if s.projects == nil {
		return
	}

	fmt.Printf("%-10s %-30s %-10s %-20s\n", "Name", "Domain", "Port", "Commands")

	for projectName, project := range s.projects {
		fmt.Printf("%-10s %-30s %-10d %-20s\n",
			projectName,
			project.Domain,
			project.Port,
			strings.Join(project.Commands, ", "),
		)
	}
}

func (s *Spinup) addCommandToProject(projectName string, commandName string) {
	exists, project := s.projectExists(projectName)

	if !exists {
		fmt.Printf("Project '%s' does not exist\n", projectName)
		return
	}

	_, exists = s.commands[commandName]

	if !exists {
		fmt.Printf("Command '%s' does not exist\n", commandName)
		return
	}

	for _, command := range project.Commands {
		if command == commandName {
			fmt.Printf("Command '%s' already exists in project '%s'\n", commandName, projectName)
			return
		}
	}

	project.Commands = append(project.Commands, commandName)

	s.projects[projectName] = project

	updatedProjectsConfig, err := json.MarshalIndent(s.projects, "", "  ")

	if err != nil {
		fmt.Println("Error encoding projects to json:", err)
		return
	}

	err = os.WriteFile(s.getProjectsFilePath(), updatedProjectsConfig, 0644)

	if err != nil {
		fmt.Println("Error writing projects to file:", err)
		return
	}

	fmt.Printf("Added command '%s' to project '%s'\n", commandName, projectName)
}

func (s *Spinup) removeCommandFromProject(projectName string, commandName string) {
	exists, project := s.projectExists(projectName)

	if !exists {
		fmt.Printf("Project '%s' does not exist\n", projectName)
		return
	}

	for i, command := range project.Commands {
		if command == commandName {
			project.Commands = append(project.Commands[:i], project.Commands[i+1:]...)

			s.projects[projectName] = project

			updatedProjectsConfig, err := json.MarshalIndent(s.projects, "", "  ")

			if err != nil {
				fmt.Println("Error encoding projects to json:", err)
				return
			}

			err = os.WriteFile(s.getProjectsFilePath(), updatedProjectsConfig, 0644)

			if err != nil {
				fmt.Println("Error writing projects to file:", err)
				return
			}

			fmt.Printf("Removed command '%s' from project '%s'\n", commandName, projectName)

			return
		}
	}
}

func (s *Spinup) setProjectDir(projectName string, dir *string) {
	exists, project := s.projectExists(projectName)

	if !exists {
		fmt.Printf("Project '%s' does not exist\n", projectName)
		return
	}

	if dir == nil {
		cwd, err := os.Getwd()

		if err != nil {
			fmt.Println("Error getting current working directory:", err)
			return
		}

		project.Dir = &cwd
	} else {
		// Check if dir exists as a directory
		info, err := os.Stat(*dir)

		if err != nil {
			fmt.Printf("Directory '%s' does not exist: %s\n", *dir, err)
			return
		}

		if !info.IsDir() {
			fmt.Printf("'%s' is not a directory\n", *dir)
			return
		}

		project.Dir = dir
	}

	s.projects[projectName] = project

	updatedProjectsConfig, err := json.MarshalIndent(s.projects, "", "  ")

	if err != nil {
		fmt.Println("Error encoding projects to json:", err)
		return
	}

	err = os.WriteFile(s.getProjectsFilePath(), updatedProjectsConfig, 0644)

	if err != nil {
		fmt.Println("Error writing projects to file:", err)
		return
	}

	fmt.Printf("Set directory to '%s' for project '%s'\n", *project.Dir, projectName)
}

func (s *Spinup) getProjectDir(projectName string) {
	exists, project := s.projectExists(projectName)

	if !exists {
		fmt.Printf("Project '%s' does not exist\n", projectName)
		return
	}

	if project.Dir == nil {
		fmt.Printf("Project '%s' does not have a directory set\n", projectName)
		return
	}

	fmt.Println(*project.Dir)
}

func (s *Spinup) handleProject() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s project <add|remove|list> [args...]\n", config.ProgramName)
		return
	}

	switch os.Args[2] {
	case "list", "ls":
		s.listProjects()
	case "add":
		if len(os.Args) == 3 {
			s.addProjectInteractive()
			return
		}

		if len(os.Args) < 6 {
			fmt.Printf("Usage: %s project add <name> <domain> <port>\n", config.ProgramName)
			return
		}

		port, err := strconv.Atoi(os.Args[5])

		if err != nil {
			fmt.Println("Port must be an integer")
			return
		}

		s.addProject(os.Args[3], os.Args[4], port, os.Args[6:])
	case "remove", "rm":
		if len(os.Args) == 3 {
			s.removeProjectInteractive()
			return
		}

		if len(os.Args) != 4 {
			fmt.Printf("Usage: %s project remove|rm <name>\n", config.ProgramName)
			return
		}

		s.removeProject(os.Args[3])
	case "add-command", "ac":
		if len(os.Args) < 5 {
			fmt.Printf("Usage: %s project add-command|ac <project> <command>\n", config.ProgramName)
			return
		}

		s.addCommandToProject(os.Args[3], os.Args[4])
	case "remove-command", "rc":
		if len(os.Args) < 5 {
			fmt.Printf("Usage: %s project remove-command|rc <project> <command>\n", config.ProgramName)
			return
		}

		s.removeCommandFromProject(os.Args[3], os.Args[4])
	case "set-dir", "sd":
		if len(os.Args) < 4 {
			fmt.Printf("Usage: %s project set-dir|sp <project> [dir]\n", config.ProgramName)
			return
		}

		if len(os.Args) == 5 {
			s.setProjectDir(os.Args[3], &os.Args[4])
			return
		}

		s.setProjectDir(os.Args[3], nil)
	case "get-dir", "gd":
		if len(os.Args) != 4 {
			fmt.Printf("Usage: %s project get-dir|gp <project>\n", config.ProgramName)
			return
		}

		s.getProjectDir(os.Args[3])
	}
}
