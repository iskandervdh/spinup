package config

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path"
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

func GetConfigDirPath() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf("could not get home directory of current user")
	}

	return path.Join(home, ".config", ProgramName), nil
}

func New() (*Config, error) {
	configDir, err := GetConfigDirPath()

	if err != nil {
		return nil, err
	}

	return &Config{
		configDir:      configDir,
		nginxConfigDir: nginxConfigDir,
		hostsFile:      hostsFile,
		hostsBackupDir: hostsBackupDir,
		testing:        false,
	}, nil
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
