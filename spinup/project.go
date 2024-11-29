package spinup

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/iskandervdh/spinup/config"
)

type Project struct {
	Domain    string    `json:"domain"`
	Port      int       `json:"port"`
	Commands  []string  `json:"commands"`
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

func (s *Spinup) addProject(name string, domain string, port int, commandNames []string) {
	if s.projects == nil {
		return
	}

	// Check if commands exist
	for _, commandName := range commandNames {
		_, exists := s.commands[commandName]

		if !exists {
			fmt.Println("Command", commandName, "does not exist")
			return
		}
	}

	// Check if project already exists or if domain or port is already in use
	for projectName, project := range s.projects {
		if projectName == name {
			fmt.Printf("Project '%s' already exists.\n", name)
			return
		}

		if project.Domain == domain {
			fmt.Printf("Project with domain '%s' already exists.\n", domain)
			return

		}

		if project.Port == port {
			fmt.Printf("Project with port '%d' already exists.\n", port)
			return
		}
	}

	err := config.AddNginxConfig(name, domain, port)

	if err != nil {
		fmt.Println("Error trying to create nginx config file", err)
		return
	}

	err = config.AddHost(domain)

	if err != nil {
		fmt.Println("Error trying to add domain to hosts file", err)

		// Remove nginx config file if adding domain to hosts file fails
		config.RemoveNginxConfig(name)

		return
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
		fmt.Println("Error encoding projects to json:", err)
		return
	}

	os.WriteFile(s.getProjectsFilePath(), updatedProjectsConfig, 0644)

	fmt.Printf("Added project '%s' with domain '%s' and port %d\n", name, domain, port)
}

func (s *Spinup) removeProject(name string) {
	if s.projects == nil {
		return
	}

	exists, _ := s.projectExists(name)

	if !exists {
		fmt.Printf("Project '%s' does not exist, nothing to remove\n", name)
		return
	}

	err := config.RemoveNginxConfig(name)

	if err != nil {
		fmt.Println("Could not remove nginx config file:", err)
	}

	err = config.RemoveHost(s.projects[name].Domain)

	if err != nil {
		fmt.Println("Error trying to remove domain from hosts file:", err)

		// Remove nginx config file if adding domain to hosts file fails
		config.RemoveNginxConfig(name)

		return
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
		fmt.Println("Error encoding projects to json:", err)
		return
	}

	os.WriteFile(s.getProjectsFilePath(), updatedProjectsConfig, 0644)

	fmt.Printf("Removed project '%s'\n", name)
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

func (s *Spinup) handleProject() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s project <add|remove|list> [args...]\n", config.ProgramName)
		return
	}

	switch os.Args[2] {
	case "list", "ls":
		s.listProjects()
	case "add":
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
		if len(os.Args) != 4 {
			fmt.Printf("Usage: %s project remove|rm <name>\n", config.ProgramName)
			return
		}

		s.removeProject(os.Args[3])
	}
}
