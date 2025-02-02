package app

import (
	"fmt"
	"os"
	"sort"

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

	// Sort projects by name
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})

	return projects
}

func (a *App) AddProject(name string, port int64, commandNames []string, projectDir string) error {
	err := a.core.FetchCommands()

	if err != nil {
		return fmt.Errorf("error getting commands config: %s", err)
	}

	err = a.core.FetchProjects()

	if err != nil {
		return fmt.Errorf("error getting projects config: %s", err)
	}

	msg := a.core.AddProject(name, port, commandNames)

	if _, ok := msg.(*common.ErrMsg); ok {
		fmt.Println(msg.GetText())
		return fmt.Errorf("%s", msg.GetText())
	}

	if projectDir != "" {
		msg = a.core.SetProjectDir(name, &projectDir)

		if _, ok := msg.(*common.ErrMsg); ok {
			fmt.Println(msg.GetText())
			return fmt.Errorf("%s", msg.GetText())
		}
	}

	return nil
}

func (a *App) UpdateProject(id int64, name string, port int64, commandNames []string, projectDir string) error {
	err := a.core.FetchCommands()

	if err != nil {
		return fmt.Errorf("error getting commands config: %s", err)
	}

	err = a.core.FetchProjects()

	if err != nil {
		return fmt.Errorf("error getting projects config: %s", err)
	}

	msg := a.core.UpdateProjectByID(id, name, port, commandNames)

	if _, ok := msg.(*common.ErrMsg); ok {
		fmt.Println(msg.GetText())
		return fmt.Errorf("%s", msg.GetText())
	}

	if projectDir != "" {
		msg = a.core.SetProjectDir(name, &projectDir)

		if _, ok := msg.(*common.ErrMsg); ok {
			fmt.Println(msg.GetText())
			return fmt.Errorf("%s", msg.GetText())
		}
	}

	return nil
}

func (a *App) RemoveProject(id int64) error {
	err := a.core.FetchCommands()

	if err != nil {
		return fmt.Errorf("error getting commands config: %s", err)
	}

	err = a.core.FetchProjects()

	if err != nil {
		return fmt.Errorf("error getting projects config: %s", err)
	}

	msg := a.core.RemoveProjectById(id)

	if _, ok := msg.(*common.ErrMsg); ok {
		fmt.Println(msg.GetText())
		return fmt.Errorf("%s", msg.GetText())
	}

	return nil
}

func (a *App) UpdateProjectDirectory(projectName string, defaultDir string) error {
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

func (a *App) SelectProjectDirectory(projectName string, defaultDir string) (string, error) {
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
		return "", err
	}

	if dir == "" {
		return "", fmt.Errorf("no directory selected")
	}

	return dir, nil
}
