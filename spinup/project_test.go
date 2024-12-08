package spinup

import (
	"os"
	"testing"

	"github.com/iskandervdh/spinup/config"
)

func TestGetProjectNamesEmpty(t *testing.T) {
	s := TestingSpinup("get_project_names_empty")

	projectNames := s.getProjectNames()

	if len(projectNames) != 0 {
		t.Error("Expected no project names, got", len(projectNames))
		return
	}
}

func TestAddProject(t *testing.T) {
	s := TestingSpinup("add_project")

	// Fetch the commands and projects from their config files
	s.getCommandsConfig()
	s.getProjectsConfig()

	s.addProject("test", "test.local", 1234, []string{})

	// "Refetch" the projects from the config file
	s.getProjectsConfig()

	projectNames := s.getProjectNames()

	if len(projectNames) != 1 {
		t.Error("Expected 1 project name, got", len(projectNames))
		return
	}

	if projectNames[0] != "test" {
		t.Error("Expected project name to be 'test', got", projectNames[0])
	}

	// Check if nginx config file was updated
	nginxFilePath := s.getConfig().GetNginxConfigDir() + "/test.conf"

	if _, err := os.Stat(nginxFilePath); os.IsNotExist(err) {
		t.Error("Expected nginx config file to exist, got", err)
		return
	}

	// Check if hosts file was updated
	hostsFilePath := s.getConfig().GetHostsFile()

	hostsFile, err := os.Open(hostsFilePath)

	if err != nil {
		t.Error("Expected to open hosts file, got", err)
		return
	}

	defer hostsFile.Close()

	fileInfo, err := os.Stat(hostsFilePath)

	if err != nil {
		t.Error("Expected to stat hosts file, got", err)
		return
	}

	buf := make([]byte, fileInfo.Size())

	_, err = hostsFile.Read(buf)

	if err != nil {
		t.Error("Expected to read hosts file, got", err)
		return
	}

	hostsContent := string(buf)
	expected := "\n\n" + config.HostsBeginMarker + "\n127.0.0.1\ttest.local" + config.HostsEndMarker

	if hostsContent != expected {
		t.Error("Expected hosts file to contain", expected, "got", hostsContent)
		return
	}
}

func TestAddProjectWithCommands(t *testing.T) {
	s := TestingSpinup("add_project_with_commands")

	// Fetch the commands and projects from their config files
	s.getCommandsConfig()
	s.getProjectsConfig()

	s.addCommand("ls", "ls")
	s.addCommand("pwd", "pwd")

	s.addProject("test", "test.local", 1234, []string{"ls", "pwd"})

	// "Refetch" the projects from the config file
	s.getProjectsConfig()

	projectNames := s.getProjectNames()

	if len(projectNames) != 1 {
		t.Error("Expected 1 project name, got", len(projectNames))
		return
	}

	if projectNames[0] != "test" {
		t.Error("Expected project name to be 'test', got", projectNames[0])
	}

	commands := s.getCommandsForProject("test")

	if len(commands) != 2 {
		t.Error("Expected 2 commands, got", len(commands))
		return
	}

	if commands[0] != "ls" {
		t.Error("Expected first command to be 'ls', got", commands[0])
	}

	if commands[1] != "pwd" {
		t.Error("Expected second command to be 'pwd', got", commands[1])
	}
}

func TestRemoveProject(t *testing.T) {
	s := TestingSpinup("remove_project")

	// Fetch the commands and projects from their config files
	s.getCommandsConfig()
	s.getProjectsConfig()

	s.addProject("test", "test.local", 1234, []string{})

	// "Refetch" the projects from the config file
	s.getProjectsConfig()

	s.removeProject("test")

	// "Refetch" the projects from the config file
	s.getProjectsConfig()

	projectNames := s.getProjectNames()

	if len(projectNames) != 0 {
		t.Error("Expected no project names, got", len(projectNames))
		return
	}

	// Check if nginx config file was removed
	nginxFilePath := s.getConfig().GetNginxConfigDir() + "/test.conf"

	if _, err := os.Stat(nginxFilePath); !os.IsNotExist(err) {
		t.Error("Expected nginx config file to not exist, got", err)
		return
	}

	// Check if hosts file was updated
	hostsFilePath := s.getConfig().GetHostsFile()
	hostsFile, err := os.Open(hostsFilePath)

	if err != nil {
		t.Error("Expected to open hosts file, got", err)
		return
	}

	defer hostsFile.Close()

	t.Error("Bogus error")
}
