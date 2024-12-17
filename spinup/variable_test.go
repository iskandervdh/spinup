package spinup

import (
	"io"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/iskandervdh/spinup/cli"
)

func TestListVariablesEmpty(t *testing.T) {
	testName := "list_variables_empty"

	r, w, err := os.Pipe()

	if err != nil {
		t.Error("Expected to create pipe, got", err)
	}

	cli := cli.New(cli.WithOut(w))
	s := TestingSpinup(testName, cli)

	s.getProjectsConfig()
	s.addProject("test", "test.local", 8000, []string{})

	// Refetch projects
	s.getProjectsConfig()

	s.listVariables("test")

	w.Close()
	out, err := io.ReadAll(r)

	if err != nil {
		t.Error("Expected to read from pipe, got", err)
	}

	if string(out) != "" {
		t.Errorf("Expected '', got '%s'", string(out))
	}
}

func TestListVariablesNoProjects(t *testing.T) {
	testName := "list_variables_no_projects"

	s := TestingSpinup(testName, nil)

	s.listVariables("test")
}

func TestListVariablesProjectNotFound(t *testing.T) {
	testName := "list_variables_project_not_found"

	s := TestingSpinup(testName, nil)

	s.getProjectsConfig()

	err := s.listVariables("test")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestListVariables(t *testing.T) {
	testName := "list_variables"

	r, w, err := os.Pipe()

	if err != nil {
		t.Error("Expected to create pipe, got", err)
	}

	cli := cli.New(cli.WithOut(w))
	s := TestingSpinup(testName, cli)

	s.getProjectsConfig()
	s.addProject("test", "test.local", 8000, []string{})
	s.addVariable("test", "key", "value")

	// Refetch projects
	s.getProjectsConfig()

	s.listVariables("test")

	w.Close()
	out, err := io.ReadAll(r)

	if err != nil {
		t.Error("Expected to read from pipe, got", err)
	}

	if strings.Contains(string(out), "value") {
		t.Errorf("Expected to get value in output, got '%s'", string(out))
	}
}

func TestAddVariable(t *testing.T) {
	testName := "add_variable"

	s := TestingSpinup(testName, nil)

	s.getProjectsConfig()
	s.addProject("test", "test.local", 8000, []string{})

	// Refetch projects
	s.getProjectsConfig()

	s.addVariable("test", "key", "value")

	projectExists, _ := s.projectExists("test")

	if !projectExists {
		t.Error("Expected project to exist, got false")
	}

	if _, ok := s.projects["test"].Variables["key"]; !ok {
		t.Error("Expected key to exist, got false")
	}
}

func TestAddVariableNoProjects(t *testing.T) {
	testName := "add_variable_no_projects"

	s := TestingSpinup(testName, nil)

	err := s.addVariable("test", "key", "value")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAddVariableProjectNotFound(t *testing.T) {
	testName := "add_variable_project_not_found"

	s := TestingSpinup(testName, nil)

	s.getProjectsConfig()

	err := s.addVariable("test", "key", "value")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAddVariableAlreadyExists(t *testing.T) {
	testName := "add_variable_already_exists"

	s := TestingSpinup(testName, nil)

	s.getProjectsConfig()
	s.addProject("test", "test.local", 8000, []string{})
	s.addVariable("test", "key", "value")

	// Refetch projects
	s.getProjectsConfig()

	err := s.addVariable("test", "key", "value")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAddVariableErrorWriting(t *testing.T) {
	testName := "add_variable_error_writing"

	s := TestingSpinup(testName, nil)

	s.getProjectsConfig()
	s.addProject("test", "test.local", 8000, []string{})

	// Refetch projects
	s.getProjectsConfig()

	// Change permissions of the projects json file to make writing fail
	os.Chmod(path.Join(s.getConfig().GetConfigDir(), "projects.json"), 0444)

	err := s.addVariable("test", "key", "value")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestRemoveVariable(t *testing.T) {
	testName := "remove_variable"

	s := TestingSpinup(testName, nil)

	s.getProjectsConfig()
	s.addProject("test", "test.local", 8000, []string{})
	s.addVariable("test", "key", "value")

	// Refetch projects
	s.getProjectsConfig()

	s.removeVariable("test", "key")

	projectExists, _ := s.projectExists("test")

	if !projectExists {
		t.Error("Expected project to exist, got false")
	}

	if _, ok := s.projects["test"].Variables["key"]; ok {
		t.Error("Expected key to not exist, got true")
	}
}

func TestRemoveVariableNoProjects(t *testing.T) {
	testName := "remove_variable_no_projects"

	s := TestingSpinup(testName, nil)

	err := s.removeVariable("test", "key")

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if err.Error() != "no projects found" {
		t.Errorf("Expected 'no projects found', got '%s'", err.Error())
	}
}
