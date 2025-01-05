package common

import (
	_ "embed"
)

const ProgramName = "spinup"
const AppCommand = "spinup-app"

//go:embed .version
var Version string
