package config

import (
	_ "embed"
	"fmt"
	"os"
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
		nginxConfigDir: getNginxConfigDir(configDir),
		testing:        false,
	}, nil
}

// Create a new Config instance with the given configuration directory.
// Used for testing purposes.
func NewTesting(testingConfigDir string) *Config {
	return &Config{
		configDir:      testingConfigDir,
		nginxConfigDir: path.Join(testingConfigDir, "/config/nginx/"),
		testing:        true,
	}
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

func (c *Config) writeToFile(filePath string, contents string) error {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		return err
	}

	_, err = f.Write([]byte(contents))

	if err != nil {
		return err
	}

	return f.Close()
}
