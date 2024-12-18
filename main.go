package main

import (
	"os"
	"os/exec"

	"github.com/iskandervdh/spinup/cli"
)

func launchApp() {
	exec.Command("spinup-app").Run()
}

func appInstalled() bool {
	_, err := exec.LookPath("spinup-app")

	return err != nil
}

func main() {
	if len(os.Args) == 1 && appInstalled() {
		launchApp()
		return
	}

	s := cli.New()
	s.Handle()
}
