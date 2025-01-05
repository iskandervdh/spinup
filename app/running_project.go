package app

import (
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"syscall"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/core"
)

type runningProject struct {
	core        *core.Core
	msgChan     *chan common.Msg
	msgChanWg   *sync.WaitGroup
	logFile     *os.File
	readingLogs bool
}

func NewRunningProject() *runningProject {
	// Create output file to write logs to
	msgChan := make(chan common.Msg, 100)
	msgChanWg := sync.WaitGroup{}
	core := core.New(core.WithMsgChan(&msgChan))

	core.FetchCommands()
	core.FetchProjects()

	rp := &runningProject{
		core:      core,
		msgChan:   &msgChan,
		msgChanWg: &msgChanWg,
	}

	rp.msgChanWg.Add(1)

	go func() {
		defer msgChanWg.Done()
		defer rp.logFile.Close()

		for msg := range *rp.msgChan {
			// If there is no log file linked to the currently running project
			// we can just forward the output to stdout
			if rp.logFile == nil {
				fmt.Println(msg.GetText())

				continue
			}

			// Write the output of the running project to the log file specified on the running project instance
			_, err := rp.logFile.Write([]byte(fmt.Sprintln(msg.GetText())))

			if err != nil {
				fmt.Println("Error writing to log file:", err)
			}
		}
	}()

	return rp
}

func (rp *runningProject) useLogFile(projectName string) error {
	configDir := rp.core.GetConfig().GetConfigDir()
	logDirPath := path.Join(configDir, "logs")

	// Make a log directory for the project if it does not exist
	err := os.MkdirAll(logDirPath, 0755)

	if err != nil {
		return err
	}

	// Determine the name of the log file
	logFilePath := path.Join(logDirPath, projectName+".log")
	logFile, err := os.Create(logFilePath)

	if err != nil {
		return err
	}

	rp.logFile = logFile

	rp.core.SetOut(rp.logFile)
	rp.core.SetErr(rp.logFile)

	return nil
}

func (rp *runningProject) GetLogFilePath() (string, error) {
	if rp.logFile == nil {
		return "", errors.New("no logs available for project")
	}

	return rp.logFile.Name(), nil
}

func (a *App) RunProject(projectName string) error {
	if a.runningProjects == nil {
		a.runningProjects = make(map[string]*runningProject)
	}

	if _, ok := a.runningProjects[projectName]; ok {
		return fmt.Errorf("project '%s' is already running", projectName)
	}

	runningProject := NewRunningProject()

	a.runningProjects[projectName] = runningProject

	err := runningProject.useLogFile(projectName)

	if err != nil {
		return err
	}

	msg := runningProject.core.TryToRun(projectName)

	if _, ok := msg.(*common.ErrMsg); ok {
		fmt.Println(msg)
		return errors.New(msg.GetText())
	} else if msg == nil {
		return fmt.Errorf("project '%s' does not exist", projectName)
	}

	return errors.New(msg.GetText())
}

func (a *App) StopProject(projectName string) string {
	if a.runningProjects == nil {
		return fmt.Sprintf("project '%s' is not running", projectName)
	}

	runningProject, ok := a.runningProjects[projectName]

	if !ok {
		return fmt.Sprintf("project '%s' is not running", projectName)
	}

	sigChan := runningProject.core.GetSigChan()
	*sigChan <- syscall.SIGINT

	if _, ok := a.runningProjects[projectName]; ok {
		delete(a.runningProjects, projectName)

		return fmt.Sprintf("project '%s' stopped", projectName)
	}

	return fmt.Sprintf("project '%s' is not running", projectName)
}
