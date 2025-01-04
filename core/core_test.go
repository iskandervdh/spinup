package core

import (
	"os"
	"path"
	"testing"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
)

func TestingConfigDir(testName string) string {
	return path.Join(os.TempDir(), common.ProgramName, testName)
}

func TestingCore(testName string) *Core {
	// Remove old tmp config dir
	testingConfigDir := TestingConfigDir(testName)
	err := os.RemoveAll(testingConfigDir)

	if err != nil {
		panic(err)
	}

	cfg := config.NewTesting(testingConfigDir)
	c := New(WithConfig(cfg))

	// Mock msgChan to prevent blocking during testing
	c.msgChan = new(chan common.Msg)
	*c.msgChan = make(chan common.Msg)

	go func() {
		for {
			<-(*c.msgChan)
		}
	}()

	// Mock init to prevent errors during testing
	c.Init()

	return c
}

func TestNew(t *testing.T) {
	testName := "new"
	c := TestingCore(testName)

	if c == nil {
		t.Error("Expected Spinup instance, got nil")
		return
	}

	if c.GetConfig() == nil {
		t.Error("Expected Config instance, got nil")
		return
	}

	if c.GetConfig().GetConfigDir() != TestingConfigDir(testName) {
		t.Error(
			"Expected ConfigDir to be",
			TestingConfigDir(testName),
			"got",
			c.GetConfig().GetConfigDir(),
		)
	}
}

func TestInit(t *testing.T) {
	c := TestingCore("init")

	configDir := c.GetConfig().GetConfigDir()

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("Expected ConfigDir to exist, got", err)
	}
}

func TestHostsConfigInit(t *testing.T) {
	c := TestingCore("hosts_config_init")

	hostsFilePath := c.GetConfig().GetHostsFile()

	if _, err := os.Stat(hostsFilePath); os.IsNotExist(err) {
		t.Error("Expected hosts file to exist, got", err)
		return
	}

	hostsFile, err := os.Open(hostsFilePath)

	if err != nil {
		t.Error("Expected to open hosts file, got", err)
		return
	}

	defer hostsFile.Close()

	// Expect hosts file to contain start and end comments
	expected := "\n\n" + config.HostsBeginMarker + config.HostsEndMarker

	buf := make([]byte, len(expected))

	_, err = hostsFile.Read(buf)

	if err != nil {
		t.Error("Expected to read hosts file, got", err)
		return
	}

	if string(buf) != expected {
		t.Error("Expected hosts file to contain", expected, "got", string(buf))
		return
	}
}

func TestHostsBackupDirInit(t *testing.T) {
	c := TestingCore("hosts_backup_dir_init")

	hostsBackupDir := c.GetConfig().GetHostsBackupDir()

	if _, err := os.Stat(hostsBackupDir); os.IsNotExist(err) {
		t.Error("Expected hosts backup dir to exist, got", err)
		return
	}
}
