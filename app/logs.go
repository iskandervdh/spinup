package app

import (
	"bufio"
	"fmt"
	"io"
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

	logChannel := make(chan string, 100)
	defer close(logChannel)

	go func() {
		for logLine := range logChannel {
			runtime.EventsEmit(a.ctx, "log", logLine)
		}
	}()

	for runningProject.readingLogs {
		line, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				time.Sleep(100 * time.Millisecond) // Reduce sleep time for faster response
				continue
			}

			return err // Return error if it's not EOF
		}

		logChannel <- line
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
