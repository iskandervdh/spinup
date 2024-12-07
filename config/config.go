package config

import _ "embed"

var ProgramName = "spinup"

//go:embed .version
var Version string

var CommandsFileName = "commands.json"

var ProjectsFileName = "projects.json"
