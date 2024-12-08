package config

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/iskandervdh/spinup/cli"
)

var ProgramName = "spinup"

//go:embed .version
var Version string

var commandsFileName = "commands.json"

var projectsFileName = "projects.json"

var hostsFile = "/etc/hosts"

var hostsBackupDir = "/etc/hosts_bak"

type Config struct {
	ConfigDir        string
	CommandsFileName string
	ProjectsFileName string
	hostsBeginMarker string
	hostsEndMarker   string
}

func getConfigDirPath() string {
	home, err := os.UserHomeDir()

	if err != nil {
		cli.ErrorPrint("Cloud not get home directory of current user")
		panic(err)
	}

	return path.Join(home, ".config", ProgramName)
}

func New(options ...func(*Config)) *Config {
	config := &Config{
		ConfigDir:        getConfigDirPath(),
		CommandsFileName: commandsFileName,
		ProjectsFileName: projectsFileName,
	}

	for _, option := range options {
		option(config)
	}

	config.hostsBeginMarker = fmt.Sprintf("### BEGIN_%s_HOSTS\n", strings.ToUpper(ProgramName))
	config.hostsEndMarker = fmt.Sprintf("\n### END_%s_HOSTS", strings.ToUpper(ProgramName))

	return config
}

func WithCommandsFileName(commandsFileName string) func(*Config) {
	return func(c *Config) {
		c.CommandsFileName = commandsFileName
	}
}

func WithProjectsFileName(projectsFileName string) func(*Config) {
	return func(c *Config) {
		c.ProjectsFileName = projectsFileName
	}
}
