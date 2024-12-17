package core

import (
	"testing"
)

func TestRun(t *testing.T) {
	s := TestingCore("run")

	s.GetCommandsConfig()
	s.GetProjectsConfig()

	s.AddCommand("ls", "ls")

	s.AddProject("test", "test.local", 1234, []string{"ls"})

	// "Refetch" the projects from the config file
	s.GetProjectsConfig()

	s.TryToRun("test")
}
