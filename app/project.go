package app

import (
	"fmt"
	"os"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/core"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) GetProjects() []core.Project {
	err := a.core.FetchCommands()
	if err != nil {
		fmt.Println("Error getting commands config:", err)
	}

	err = a.core.FetchProjects()

	if err != nil {
		fmt.Println("Error getting projects config:", err)
	}

	projects, err := a.core.GetProjects()

	if err != nil {
		fmt.Println("Error getting projects:", err)

		return nil
	}

	return projects
}

func (a *App) AddProject(name string, domain string, port int64, commandNames []string) error {
	err := a.core.FetchCommands()

	if err != nil {
		return fmt.Errorf("error getting commands config: %s", err)
	}

	err = a.core.FetchProjects()

	if err != nil {
		return fmt.Errorf("error getting projects config: %s", err)
	}

	msg := a.core.AddProject(name, domain, port, commandNames)

	if _, ok := msg.(*common.ErrMsg); ok {
		fmt.Println(msg.GetText())
		return fmt.Errorf("%s", msg.GetText())
	}

	return nil
}

func (a *App) UpdateProject(name string, domain string, port int64, commandNames []string) error {
	err := a.core.FetchCommands()

	if err != nil {
		return fmt.Errorf("error getting commands config: %s", err)
	}

	err = a.core.FetchProjects()

	if err != nil {
		return fmt.Errorf("error getting projects config: %s", err)
	}

	msg := a.core.UpdateProject(name, domain, port, commandNames)

	if _, ok := msg.(*common.ErrMsg); ok {
		fmt.Println(msg.GetText())
		return fmt.Errorf("%s", msg.GetText())
	}

	return nil
}

func (a *App) RemoveProject(name string) error {
	err := a.core.FetchCommands()

	if err != nil {
		return fmt.Errorf("error getting commands config: %s", err)
	}

	err = a.core.FetchProjects()

	if err != nil {
		return fmt.Errorf("error getting projects config: %s", err)
	}

	msg := a.core.RemoveProject(name)

	if _, ok := msg.(*common.ErrMsg); ok {
		fmt.Println(msg.GetText())
		return fmt.Errorf("%s", msg.GetText())
	}

	return nil
}

func (a *App) SelectProjectDirectory(projectName string, defaultDir string) error {
	if defaultDir == "" {
		d, err := os.UserHomeDir()

		if err != nil {
			defaultDir = ""
		} else {
			defaultDir = d
		}
	}

	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title:            "Select a directory for project " + projectName,
		DefaultDirectory: defaultDir,
	})

	if err != nil {
		fmt.Println("Error selecting directory:", err)
		return err
	}

	if dir == "" {
		return fmt.Errorf("no directory selected")
	}

	msg := a.core.SetProjectDir(projectName, &dir)

	if _, ok := msg.(*common.ErrMsg); ok {
		fmt.Println(msg.GetText())
		return fmt.Errorf("%s", msg.GetText())
	}

	return nil
}
