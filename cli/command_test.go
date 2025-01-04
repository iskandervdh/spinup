package cli

import (
	"os"
	"testing"

	"github.com/iskandervdh/spinup/config"
)

func TestHandleCommandTooFewArguments(t *testing.T) {
	c := New()

	os.Args = []string{config.ProgramName, "command"}

	c.Handle()
}

func TestHandleCommandLs(t *testing.T) {
	c := New()

	os.Args = []string{config.ProgramName, "command", "list"}
	c.Handle()

	c = New()

	os.Args = []string{config.ProgramName, "c", "ls"}
	c.Handle()
}
