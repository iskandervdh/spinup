package app

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/core"
)

type App struct {
	ctx             context.Context
	core            *core.Core
	runningProjects map[string]*runningProject
}

func NewApp() *App {
	return &App{}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.core = core.New()

	err := a.core.FetchCommands()

	if err != nil {
		fmt.Println("Error getting commands config:", err)
	}

	err = a.core.FetchProjects()

	if err != nil {
		fmt.Println("Error getting projects config:", err)
	}
}

func (a *App) GetSpinupVersion() string {
	return common.Version
}
