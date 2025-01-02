package core

import (
	"testing"
)

func TestRun(t *testing.T) {
	c := TestingCore("run")

	c.FetchCommands()
	c.FetchProjects()

	c.AddCommand("ls", "ls")

	c.AddProject("test", "test.local", 1234, []string{"ls"})

	// "Refetch" the projects from the config file
	c.FetchProjects()

	c.TryToRun("test")
}
