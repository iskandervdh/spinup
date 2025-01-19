package common

import (
	"runtime"
)

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

func GetDomain(projectName string) string {
	return projectName + "." + TLD
}
