package common

import (
	_ "embed"
)

var ProgramName = "spinup"
var AppCommand = "spinup-app"

//go:embed .version
var Version string
