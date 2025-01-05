package common

import (
	_ "embed"
)

const ProgramName = "spinup"

//go:embed .version
var Version string
