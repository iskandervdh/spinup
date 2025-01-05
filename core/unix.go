//go:build linux || darwin

package core

import (
	"os"
	"syscall"
)

func createProcessGroup() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}

func killProcess(process *os.Process) error {
	return syscall.Kill(-process.Pid, syscall.SIGTERM)
}
