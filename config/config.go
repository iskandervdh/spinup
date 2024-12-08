package config

import (
	_ "embed"
	"os"
	"os/exec"
	"path"

	"github.com/iskandervdh/spinup/cli"
)

var ProgramName = "spinup"

//go:embed .version
var Version string

var CommandsFileName = "commands.json"

var ProjectsFileName = "projects.json"

var nginxConfigDir = "/etc/nginx/conf.d"

var hostsFile = "/etc/hosts"

var hostsBackupDir = "/etc/hosts_bak"

type Config struct {
	configDir      string
	nginxConfigDir string
	hostsFile      string
	hostsBackupDir string

	testing bool
}

func GetConfigDirPath() string {
	home, err := os.UserHomeDir()

	if err != nil {
		cli.ErrorPrint("Cloud not get home directory of current user")
		panic(err)
	}

	return path.Join(home, ".config", ProgramName)
}

func New() *Config {
	return &Config{
		configDir:      GetConfigDirPath(),
		nginxConfigDir: nginxConfigDir,
		hostsFile:      hostsFile,
		hostsBackupDir: hostsBackupDir,
		testing:        false,
	}
}

func NewTesting(testingConfigDir string) *Config {
	return &Config{
		configDir:      testingConfigDir,
		nginxConfigDir: path.Join(testingConfigDir, "/nginx/conf.d"),
		hostsFile:      path.Join(testingConfigDir, "hosts"),
		hostsBackupDir: path.Join(testingConfigDir, "hosts_bak"),
		testing:        true,
	}
}

func (c *Config) withSudo(name string, args ...string) *exec.Cmd {
	if c.IsTesting() {
		return exec.Command(name, args...)
	}

	return exec.Command("sudo", append([]string{name}, args...)...)
}

/**
 * Getters
 */
func (c *Config) GetConfigDir() string {
	return c.configDir
}

func (c *Config) GetNginxConfigDir() string {
	return c.nginxConfigDir
}

func (c *Config) GetHostsFile() string {
	return c.hostsFile
}

func (c *Config) GetHostsBackupDir() string {
	return c.hostsBackupDir
}

func (c *Config) IsTesting() bool {
	return c.testing
}
