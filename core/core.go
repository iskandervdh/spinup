package core

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
	"github.com/iskandervdh/spinup/database"
	"github.com/iskandervdh/spinup/database/sqlc"
	_ "github.com/mattn/go-sqlite3"
)

// Core is the main struct for the core package.
//
// It contains the config, message channel, commands and projects.
// The config is used to read and update the configuration.
// The message channel is used to send messages to the CLI or the app.
type Core struct {
	config  *config.Config
	msgChan *chan common.Msg
	sigChan *chan os.Signal

	out io.Writer
	err io.Writer

	dbQueries *sqlc.Queries
	dbContext context.Context

	commands Commands
	projects Projects
}

func (c *Core) connectToDB(config *config.Config) (*sql.DB, error) {
	databasePath := config.GetDatabasePath()

	_, err := os.Stat(databasePath)

	if err != nil {
		if os.IsNotExist(err) {
			// fmt.Println("Database does not exist, initializing...")
			err := c.InitSqliteDB()

			if err != nil {
				return nil, fmt.Errorf("error initializing database: %s", err)
			}
		} else {
			return nil, fmt.Errorf("error getting database info: %s", err)
		}
	}

	db, err := sql.Open("sqlite3", databasePath)

	if err != nil {
		return nil, fmt.Errorf("error opening database: %s", err)
	}

	return db, nil
}

// Create a new Core instance with the given options.
func New(options ...func(*Core)) *Core {
	config, err := config.New()

	if err != nil {
		fmt.Println("Error getting config:", err)
		os.Exit(1)
	}

	c := &Core{
		config: config,
		out:    os.Stdout,
		err:    os.Stderr,

		dbContext: context.Background(),

		commands: nil,
		projects: nil,
	}

	for _, option := range options {
		option(c)
	}

	db, err := c.connectToDB(c.config)

	if err != nil {
		fmt.Println("Error connecting to database:", err)
		os.Exit(1)
	}

	c.dbQueries = sqlc.New(db)
	err = database.MigrateDatabase(db)

	if err != nil {
		fmt.Println("Error migrating database:", err)
		os.Exit(1)
	}

	return c
}

// Optional function to set the config of the Core when creating a new instance.
func WithConfig(config *config.Config) func(*Core) {
	return func(c *Core) {
		c.config = config
	}
}

// Optional function to set the message channel of the Core when creating a new instance.
func WithMsgChan(msgChan *chan common.Msg) func(*Core) {
	return func(c *Core) {
		c.msgChan = msgChan
	}
}

func (c *Core) SetOut(out io.Writer) {
	c.out = out
}

func (c *Core) SetErr(err io.Writer) {
	c.err = err
}

// Get the config of the Core instance.
func (c *Core) GetConfig() *config.Config {
	return c.config
}

// RequireSudo checks if the user has sudo permissions.
// This is needed to update some system files.
//
// It returns an error if the user does not have sudo permissions.
func (c *Core) RequireSudo() error {
	if c.config.IsTesting() {
		return nil
	}

	err := exec.Command("sudo", "-v").Run()

	if err != nil {
		return fmt.Errorf("this command requires sudo")
	}

	return nil
}

// Get all the names of the commands.
func (c *Core) GetCommandNames() []string {
	if c.commands == nil {
		c.FetchCommands()
	}

	var commandNames []string

	for _, command := range c.commands {
		commandNames = append(commandNames, command.Name)
	}

	return commandNames
}

// Get all the names of the projects.
func (c *Core) GetProjectNames() []string {
	if c.projects == nil {
		c.FetchProjects()
	}

	var projectNames []string

	for _, project := range c.projects {
		projectNames = append(projectNames, project.Name)
	}

	return projectNames
}

// Get the signal channel of the Core instance.
func (c *Core) GetSigChan() *chan os.Signal {
	return c.sigChan
}

// Get all the commands that are part of the given project.
func (c *Core) getCommandsForProject(projectName string) []Command {
	if c.projects == nil {
		c.FetchProjects()
	}

	index := slices.IndexFunc(c.projects, func(p Project) bool {
		return p.Name == projectName
	})

	if index == -1 {
		return nil
	}

	project := c.projects[index]

	return project.Commands
}

// Send a message to the message channel.
func (c *Core) sendMsg(msg common.Msg) {
	*c.msgChan <- msg
}
