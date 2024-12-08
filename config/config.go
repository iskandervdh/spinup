package config

import (
	_ "embed"

	"github.com/iskandervdh/spinup/cli"
)

var ProgramName = "spinup"

//go:embed .version
var Version string

var CommandsFileName = "commands.json"

var ProjectsFileName = "projects.json"
