package core

import (
	"os"
	"path"
	"testing"
)

func TestListVariablesEmpty(t *testing.T) {
	testName := "list_variables_empty"

	c := TestingCore(testName)

	c.FetchProjects()
	c.AddProject("test", "test.local", 8000, []string{})

	// Refetch projects
	c.FetchProjects()

	exists, project := c.ProjectExists("test")

	if !exists {
		t.Error("Expected project to exist, got false")
	}

	variables := project.Variables

	if !(len(variables) == 0) {
		t.Error("Expected variables to be empty, got false")
	}
}

// func TestListVariablesProjectNotFound(t *testing.T) {
// 	testName := "list_variables_project_not_found"

// 	s := TestingCore(testName)

// 	s.FetchProjects()

// 	err := s.listVariables("test")

// 	if err == nil {
// 		t.Error("Expected error, got nil")
// 	}
// }

// func TestListVariables(t *testing.T) {
// 	testName := "list_variables"

// 	s := TestingCore(testName)

// 	s.FetchProjects()
// 	s.AddProject("test", "test.local", 8000, []string{})
// 	s.AddVariable("test", "key", "value")

// 	// Refetch projects
// 	s.FetchProjects()

// 	s.listVariables("test")
// }

func TestAddVariable(t *testing.T) {
	testName := "add_variable"

	s := TestingCore(testName)

	s.FetchProjects()
	s.AddProject("test", "test.local", 8000, []string{})

	// Refetch projects
	s.FetchProjects()

	s.AddVariable("test", "key", "value")

	s.FetchProjects()

	exists, project := s.ProjectExists("test")

	if !exists {
		t.Error("Expected project to exist, got false")
	}

	if len(project.Variables) == 0 {
		t.Error("Expected variables to exist, got false")
	}

	variable := project.Variables[0]

	if variable.Value != "value" {
		t.Error("Expected key to exist, got false")
	}
}

func TestAddVariableNoProjects(t *testing.T) {
	testName := "add_variable_no_projects"

	s := TestingCore(testName)

	err := s.AddVariable("test", "key", "value")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAddVariableProjectNotFound(t *testing.T) {
	testName := "add_variable_project_not_found"

	s := TestingCore(testName)

	s.FetchProjects()

	err := s.AddVariable("test", "key", "value")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAddVariableAlreadyExists(t *testing.T) {
	testName := "add_variable_already_exists"

	s := TestingCore(testName)

	s.FetchProjects()
	s.AddProject("test", "test.local", 8000, []string{})
	s.AddVariable("test", "key", "value")

	// Refetch projects
	s.FetchProjects()

	err := s.AddVariable("test", "key", "value")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAddVariableErrorWriting(t *testing.T) {
	testName := "add_variable_error_writing"

	s := TestingCore(testName)

	s.FetchProjects()
	s.AddProject("test", "test.local", 8000, []string{})

	// Refetch projects
	s.FetchProjects()

	// Change permissions of the projects json file to make writing fail
	os.Chmod(path.Join(s.GetConfig().GetConfigDir(), "projects.json"), 0444)

	err := s.AddVariable("test", "key", "value")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestRemoveVariable(t *testing.T) {
	testName := "remove_variable"

	s := TestingCore(testName)

	s.FetchProjects()
	s.AddProject("test", "test.local", 8000, []string{})
	s.AddVariable("test", "key", "value")

	// Refetch projects
	s.FetchProjects()

	s.RemoveVariable("test", "key")

	exists, project := s.ProjectExists("test")

	if !exists {
		t.Error("Expected project to exist, got false")
	}

	if len(project.Variables) != 0 {
		t.Errorf("Expected no variable to not exist, got %d", len(project.Variables))
	}
}

func TestRemoveVariableNoProjects(t *testing.T) {
	testName := "remove_variable_no_projects"

	s := TestingCore(testName)

	err := s.RemoveVariable("test", "key")

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.GetText() != "no projects found" {
		t.Errorf("Expected 'no projects found', got '%s'", err.GetText())
	}
}
