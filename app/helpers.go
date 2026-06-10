package app

import (
	"fmt"
	"os/exec"
	"runtime"
)

func (a *App) OpenFolder(path string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", path).Start()
	case "windows":
		return exec.Command("explorer", path).Start()
	case "darwin":
		return exec.Command("open", path).Start()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}
