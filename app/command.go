package app

import (
	"fmt"
	"sort"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/core"
)

func (a *App) GetCommands() []core.Command {
	err := a.core.FetchCommands()

	if err != nil {
		fmt.Println("Error getting commands config:", err)

		return nil
	}

	commands, err := a.core.GetCommands()

	if err != nil {
		fmt.Println("Error getting commands:", err)

		return nil
	}

	// Sort commands by name
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})

	return commands
}

func (a *App) AddCommand(name string, command string) error {
	err := a.core.FetchCommands()

	if err != nil {
		return fmt.Errorf("error getting commands config: %s", err)
	}

	msg := a.core.AddCommand(name, command)

	if _, ok := msg.(*common.ErrMsg); ok {
		fmt.Println(msg.GetText())
		return fmt.Errorf("%s", msg.GetText())
	}

	return nil
}

func (a *App) UpdateCommand(id int64, name string, command string) error {
	err := a.core.FetchCommands()

	if err != nil {
		return fmt.Errorf("error getting commands config: %s", err)
	}

	msg := a.core.UpdateCommandById(id, name, command)

	if _, ok := msg.(*common.ErrMsg); ok {
		fmt.Println(msg.GetText())
		return fmt.Errorf("%s", msg.GetText())
	}

	return nil
}

func (a *App) RemoveCommand(id int64) error {
	err := a.core.FetchCommands()

	if err != nil {
		return fmt.Errorf("error getting commands config: %s", err)
	}

	msg := a.core.RemoveCommandById(id)

	if _, ok := msg.(*common.ErrMsg); ok {
		fmt.Println(msg.GetText())
		return fmt.Errorf("%s", msg.GetText())
	}

	return nil
}
