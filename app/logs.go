package app

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type LogEvent struct {
	NewLogData string `json:"newLogData"`
}

func (a *App) FollowProjectLogs(projectName string) error {
	runningProject, ok := a.runningProjects[projectName]

	if !ok {
		return fmt.Errorf("project '%s' is not running", projectName)
	}

	logFilePath, err := runningProject.GetLogFilePath()

	if err != nil {
		return err
	}

	logFile, err := os.Open(logFilePath)

	if err != nil {
		return err
	}

	defer logFile.Close()

	reader := bufio.NewReader(logFile)
	runningProject.readingLogs = true

	for runningProject.readingLogs {
		line, err := reader.ReadString('\n')

		if err != nil {
			// If EOF is reached, sleep and retry
			time.Sleep(500 * time.Millisecond)
			continue
		}

		runtime.EventsEmit(a.ctx, "log", line)
	}

	return nil
}

func (a *App) StopFollowingProjectLogs(projectName string) error {
	runningProject, ok := a.runningProjects[projectName]

	if !ok {
		return fmt.Errorf("project '%s' is not running", projectName)
	}

	runningProject.readingLogs = false

	return nil
}
