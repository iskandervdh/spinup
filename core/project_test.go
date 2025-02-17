package core

import (
	"fmt"
	"os"
	"testing"
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

	c.AddProject("test", 1234, []string{})

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
}

func TestAddProjectWithCommands(t *testing.T) {
	c := TestingCore("add_project_with_commands")

	// Fetch the commands and projects from their config files
	c.FetchCommands()
	c.FetchProjects()

	c.AddCommand("ls", "ls")
	c.AddCommand("pwd", "pwd")

	c.AddProject("test", 1234, []string{"ls", "pwd"})

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

func TestRemoveCommandFromProject(t *testing.T) {
	c := TestingCore("remove_command_from_project")

	// Fetch the commands and projects from their config files
	c.FetchCommands()
	c.FetchProjects()

	c.AddCommand("ls", "ls")
	c.AddCommand("pwd", "pwd")

	c.AddProject("test", 1234, []string{"ls", "pwd"})

	// "Refetch" the projects from the config file
	c.FetchProjects()

	c.RemoveCommandFromProject("test", "ls")

	// "Refetch" the projects from the config file
	c.FetchProjects()

	commands := c.getCommandsForProject("test")

	if len(commands) != 1 {
		t.Error("Expected 1 command, got", len(commands))
		return
	}

	if commands[0].Command != "pwd" {
		t.Error("Expected command to be 'pwd', got", commands[0])
	}
}

func TestRemoveProject(t *testing.T) {
	c := TestingCore("remove_project")

	// Fetch the commands and projects from their config files
	c.FetchCommands()
	c.FetchProjects()

	c.AddProject("test", 1234, []string{})

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
}

func TestEditProject(t *testing.T) {
	c := TestingCore("edit_project")

	// Fetch the commands and projects from their config files
	c.FetchCommands()
	c.FetchProjects()

	c.AddProject("test", 1234, []string{})

	// "Refetch" the projects from the config file
	c.FetchProjects()

	c.UpdateProject("test", 1235, []string{})

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
}

func TestRenameProject(t *testing.T) {
	c := TestingCore("rename_project")

	// Fetch the commands and projects from their config files
	c.FetchCommands()
	c.FetchProjects()

	c.AddProject("test", 1234, []string{})

	// "Refetch" the projects from the config file
	c.FetchProjects()

	msg := c.RenameProject("test", "example")

	fmt.Println(msg)

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
}
