package common

import (
	"os/exec"
	"runtime"

	"github.com/iskandervdh/spinup/config"
)

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func AppInstalled() bool {
	_, err := exec.LookPath(AppCommand)

	return err == nil
}
