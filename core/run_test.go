package core

import (
	"testing"
)

func TestRun(t *testing.T) {
	c := TestingCore("run")

	c.GetCommandsConfig()
	c.GetProjectsConfig()

	c.AddCommand("ls", "ls")

	c.AddProject("test", "test.local", 1234, []string{"ls"})

	// "Refetch" the projects from the config file
	c.GetProjectsConfig()

	c.TryToRun("test")
}
