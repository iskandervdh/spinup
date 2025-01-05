package app

import (
	"github.com/iskandervdh/spinup/core"
)

func (a *App) ProjectPlaceholder() core.Project {
	return core.Project{}
}

func (a *App) CommandPlaceholder() core.Command {
	return core.Command{}
}

func (a *App) VariablePlaceholder() core.Variable {
	return core.Variable{}
}

func (a *App) DomainAliasPlaceholder() core.DomainAlias {
	return core.DomainAlias{}
}
