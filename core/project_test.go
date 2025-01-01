package core

import (
	"os"
	"testing"

	"github.com/iskandervdh/spinup/config"
)

func TestGetProjectNamesEmpty(t *testing.T) {
	c := TestingCore("get_project_names_empty")

	projectNames := c.GetProjectNames()

	if len(projectNames) != 0 {
		t.Error("Expected no project names, got", len(projectNames))
		return
	}
}

func TestAddProject(t *testing.T) {
	c := TestingCore("add_project")

	// Fetch the commands and projects from their config files
	c.FetchCommands()
	c.FetchProjects()

	c.AddProject("test", "test.local", 1234, []string{})

	// "Refetch" the projects from the config file
	c.FetchProjects()

	projectNames := c.GetProjectNames()

	if len(projectNames) != 1 {
		t.Error("Expected 1 project name, got", len(projectNames))
		return
	}

	if projectNames[0] != "test" {
		t.Error("Expected project name to be 'test', got", projectNames[0])
	}

	// Check if nginx config file was updated
	nginxFilePath := c.GetConfig().GetNginxConfigDir() + "/test.conf"

	if _, err := os.Stat(nginxFilePath); os.IsNotExist(err) {
		t.Error("Expected nginx config file to exist, got", err)
		return
	}

	// Check if hosts file was updated
	hostsFilePath := c.GetConfig().GetHostsFile()

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
	c := TestingCore("add_project_with_commands")

	// Fetch the commands and projects from their config files
	c.FetchCommands()
	c.FetchProjects()

	c.AddCommand("ls", "ls")
	c.AddCommand("pwd", "pwd")

	c.AddProject("test", "test.local", 1234, []string{"ls", "pwd"})

	// "Refetch" the projects from the config file
	c.FetchProjects()

	projectNames := c.GetProjectNames()

	if len(projectNames) != 1 {
		t.Error("Expected 1 project name, got", len(projectNames))
		return
	}

	if projectNames[0] != "test" {
		t.Error("Expected project name to be 'test', got", projectNames[0])
	}

	commands := c.getCommandsForProject("test")

	if len(commands) != 2 {
		t.Error("Expected 2 commands, got", len(commands))
		return
	}

	if commands[0].Command != "ls" {
		t.Error("Expected first command to be 'ls', got", commands[0])
	}

	if commands[1].Command != "pwd" {
		t.Error("Expected second command to be 'pwd', got", commands[1])
	}
}

func TestRemoveProject(t *testing.T) {
	c := TestingCore("remove_project")

	// Fetch the commands and projects from their config files
	c.FetchCommands()
	c.FetchProjects()

	c.AddProject("test", "test.local", 1234, []string{})

	// "Refetch" the projects from the config file
	c.FetchProjects()

	c.RemoveProject("test")

	// "Refetch" the projects from the config file
	c.FetchProjects()

	projectNames := c.GetProjectNames()

	if len(projectNames) != 0 {
		t.Error("Expected no project names, got", len(projectNames))
		return
	}

	// Check if nginx config file was removed
	nginxFilePath := c.GetConfig().GetNginxConfigDir() + "/test.conf"

	if _, err := os.Stat(nginxFilePath); !os.IsNotExist(err) {
		t.Error("Expected nginx config file to not exist, got", err)
		return
	}

	// Check if hosts file was updated
	hostsFilePath := c.GetConfig().GetHostsFile()
	hostsFile, err := os.Open(hostsFilePath)

	if err != nil {
		t.Error("Expected to open hosts file, got", err)
		return
	}

	defer hostsFile.Close()
}

func TestEditProject(t *testing.T) {
	c := TestingCore("edit_project")

	// Fetch the commands and projects from their config files
	c.FetchCommands()
	c.FetchProjects()

	c.AddProject("test", "test.local", 1234, []string{})

	// "Refetch" the projects from the config file
	c.FetchProjects()

	c.UpdateProject("test", "example.local", 1235, []string{})

	// "Refetch" the projects from the config file
	c.FetchProjects()

	projectNames := c.GetProjectNames()

	if len(projectNames) != 1 {
		t.Error("Expected 1 project name, got", len(projectNames))
		return
	}

	// Check if nginx config file was updated
	nginxFilePath := c.GetConfig().GetNginxConfigDir() + "/test.conf"

	if _, err := os.Stat(nginxFilePath); os.IsNotExist(err) {
		t.Error("Expected nginx config file to exist, got", err)
		return
	}

	// Check if hosts file was updated
	hostsFilePath := c.GetConfig().GetHostsFile()

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
	expected := "\n\n" + config.HostsBeginMarker + "\n127.0.0.1\texample.local" + config.HostsEndMarker

	if hostsContent != expected {
		t.Error("Expected hosts file to contain", expected, "got", hostsContent)
		return
	}
}

func TestRenameProject(t *testing.T) {
	c := TestingCore("rename_project")

	// Fetch the commands and projects from their config files
	c.FetchCommands()
	c.FetchProjects()

	c.AddProject("test", "test.local", 1234, []string{})

	// "Refetch" the projects from the config file
	c.FetchProjects()

	c.RenameProject("test", "example")

	// "Refetch" the projects from the config file
	c.FetchProjects()

	projectNames := c.GetProjectNames()

	if len(projectNames) != 1 {
		t.Error("Expected 1 project name, got", len(projectNames))
		return
	}

	if projectNames[0] != "example" {
		t.Error("Expected project name to be 'example', got", projectNames[0])
		return
	}

	// Check if nginx config file was updated
	nginxFilePath := c.GetConfig().GetNginxConfigDir() + "/example.conf"

	if _, err := os.Stat(nginxFilePath); os.IsNotExist(err) {
		t.Error("Expected nginx config file to exist, got", err)
		return
	}

	// Check if hosts file was updated
	hostsFilePath := c.GetConfig().GetHostsFile()

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
