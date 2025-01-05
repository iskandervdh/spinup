//go:build windows

package core

import (
	"os"
	"syscall"
)

func createProcessGroup() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}

func killProcess(process *os.Process) error {
	return process.Kill()
}
