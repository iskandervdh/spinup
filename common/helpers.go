package common

import (
	"os/exec"
	"runtime"
)

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

func AppInstalled() bool {
	_, err := exec.LookPath(AppCommand)

	return err == nil
}
