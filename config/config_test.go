package config

import (
	"os"
	"path"
	"testing"

	"github.com/iskandervdh/spinup/common"
)

func TestingConfigDir(testName string) string {
	return path.Join(os.TempDir(), common.ProgramName, testName)
}

func TestingConfig(testName string) *Config {
	testingConfigDir := TestingConfigDir(testName)
	err := os.RemoveAll(testingConfigDir)

	if err != nil {
		panic(err)
	}

	os.MkdirAll(testingConfigDir, 0755)

	return NewTesting(testingConfigDir)
}

func TestGetConfigDirPath(t *testing.T) {
	c := TestingConfig("config_dir_path")
	path := c.GetConfigDir()

	if path == "" {
		t.Error("Expected config dir path, got empty string")
	}
}

func TestNew(t *testing.T) {
	_, err := New()

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestNewTesting(t *testing.T) {
	c := TestingConfig("testing")

	if c == nil {
		t.Error("Expected Config instance, got nil")
	}
}

func TestGetters(t *testing.T) {
	c := TestingConfig("config_getters")

	if c.GetConfigDir() == "" {
		t.Error("Expected config dir, got empty string")
	}

	if c.GetNginxConfigDir() == "" {
		t.Error("Expected nginx config dir, got empty string")
	}

	if !c.IsTesting() {
		t.Error("Expected testing to be false, got true")
	}
}
