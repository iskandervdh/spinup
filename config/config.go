package config

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/iskandervdh/spinup/common"
)

// Config contains all configuration options for the application.
type Config struct {
	configDir      string
	nginxConfigDir string

	testing bool
}

// Returns the path to the configuration directory.
func GetDefaultConfigDirPath() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf("could not get home directory of current user")
	}

	if common.IsWindows() {
		return path.Join(home, "AppData", "Roaming", common.ProgramName), nil
	}

	return path.Join(home, ".config", common.ProgramName), nil
}

// Create a new Config instance with the default configuration.
func New() (*Config, error) {
	configDir, err := GetDefaultConfigDirPath()

	if err != nil {
		return nil, err
	}

	return &Config{
		configDir:      configDir,
		nginxConfigDir: nginxConfigDir,
		testing:        false,
	}, nil
}

// Create a new Config instance with the given configuration directory.
// Used for testing purposes.
func NewTesting(testingConfigDir string) *Config {
	return &Config{
		configDir:      testingConfigDir,
		nginxConfigDir: path.Join(testingConfigDir, "/nginx/conf.d"),
		testing:        true,
	}
}

// Add sudo to a command if not in testing mode.
func (c *Config) withSudo(name string, args ...string) *exec.Cmd {
	if c.IsTesting() || common.IsWindows() {
		return exec.Command(name, args...)
	}

	return exec.Command("sudo", append([]string{name}, args...)...)
}

// Returns the path to the configuration directory.
func (c *Config) GetConfigDir() string {
	return c.configDir
}

// Returns the path of the sqlite3 database file.
func (c *Config) GetDatabasePath() string {
	return path.Join(c.configDir, common.ProgramName+".sqlite3")
}

// Returns the path to the nginx configuration directory.
func (c *Config) GetNginxConfigDir() string {
	return c.nginxConfigDir
}

// Returns whether the application is in testing mode.
func (c *Config) IsTesting() bool {
	return c.testing
}
