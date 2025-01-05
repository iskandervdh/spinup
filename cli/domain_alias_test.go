package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/iskandervdh/spinup/common"
)

func TestDomainAlias(t *testing.T) {
	output := &bytes.Buffer{}
	c := TestingCLI("domain_alias", WithOut(output), WithErr(output))

	os.Args = []string{common.ProgramName, "domain-alias"}
	c.Handle()

	// if strings.Contains(output.String(), "Usage: spinup domain-alias|da") == false {
	// 	t.Errorf("Expected usage message, got '%s'", output.String())
	// }

	c = TestingCLI("domain_alias")

	os.Args = []string{common.ProgramName, "da"}
	c.Handle()
}
