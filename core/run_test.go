package core

import (
	"testing"
)

func TestRun(t *testing.T) {
	s := TestingSpinup("run", nil)

	s.getCommandsConfig()
	s.getProjectsConfig()

	s.addCommand("ls", "ls")

	s.addProject("test", "test.local", 1234, []string{"ls"})

	// "Refetch" the projects from the config file
	s.getProjectsConfig()

	s.tryToRun("test")
}
