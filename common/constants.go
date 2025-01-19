package common

import (
	_ "embed"
)

const ProgramName = "spinup"

const TLD = "test"

//go:embed .version
var Version string
