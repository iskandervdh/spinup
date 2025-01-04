package cli

import (
	"os"
	"testing"

	"github.com/iskandervdh/spinup/common"
)

func TestHandleCommandTooFewArguments(t *testing.T) {
	c := New()

	os.Args = []string{common.ProgramName, "command"}

	c.Handle()
}

func TestHandleCommandLs(t *testing.T) {
	c := New()

	os.Args = []string{common.ProgramName, "command", "list"}
	c.Handle()

	c = New()

	os.Args = []string{common.ProgramName, "c", "ls"}
	c.Handle()
}
