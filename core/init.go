package core

import (
	"os"

	"github.com/iskandervdh/spinup/common"
)

// Create the config directory if it doesn't exist.
func (c *Core) createConfigDir() error {
	err := os.MkdirAll(c.config.GetConfigDir(), 0755)

	if err != nil {
		return err
	}

	return nil
}

func (c *Core) InitSqliteDB() error {
	c.createConfigDir()

	_, err := os.Create(c.config.GetDatabasePath())

	if err != nil {
		return err
	}

	return nil
}

// Initialize the config directory, hosts and nginx.
func (c *Core) Init() common.Msg {
	err := c.createConfigDir()

	if err != nil {
		return common.NewErrMsg("Error creating config directory: %v", err)
	}

	err = c.config.InitNginx()

	if err != nil {
		return common.NewErrMsg("Error initializing nginx: %v", err)
	}

	return common.NewSuccessMsg("\nInitialization complete")
}
