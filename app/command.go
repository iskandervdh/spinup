package app

import (
	"fmt"

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

func (a *App) UpdateCommand(name string, command string) error {
	err := a.core.FetchCommands()

	if err != nil {
		return fmt.Errorf("error getting commands config: %s", err)
	}

	msg := a.core.UpdateCommand(name, command)

	if _, ok := msg.(*common.ErrMsg); ok {
		fmt.Println(msg.GetText())
		return fmt.Errorf("%s", msg.GetText())
	}

	return nil
}

func (a *App) RemoveCommand(name string) error {
	err := a.core.FetchCommands()

	if err != nil {
		return fmt.Errorf("error getting commands config: %s", err)
	}

	msg := a.core.RemoveCommand(name)

	if _, ok := msg.(*common.ErrMsg); ok {
		fmt.Println(msg.GetText())
		return fmt.Errorf("%s", msg.GetText())
	}

	return nil
}
