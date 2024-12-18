package core

import (
	"os"
	"path"
	"testing"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
)

func TestingConfigDir(testName string) string {
	return path.Join(os.TempDir(), config.ProgramName, testName)
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
	s := TestingCore(testName)

	if s == nil {
		t.Error("Expected Spinup instance, got nil")
		return
	}

	if s.getConfig() == nil {
		t.Error("Expected Config instance, got nil")
		return
	}

	if s.getConfig().GetConfigDir() != TestingConfigDir(testName) {
		t.Error(
			"Expected ConfigDir to be",
			TestingConfigDir(testName),
			"got",
			s.getConfig().GetConfigDir(),
		)
	}
}

func TestInit(t *testing.T) {
	s := TestingCore("init")

	configDir := s.getConfig().GetConfigDir()

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("Expected ConfigDir to exist, got", err)
	}
}

func TestCommandsConfigInit(t *testing.T) {
	s := TestingCore("commands_config_init")

	commandsFilePath := path.Join(s.getConfig().GetConfigDir(), config.CommandsFileName)

	if _, err := os.Stat(commandsFilePath); os.IsNotExist(err) {
		t.Error("Expected commands config file to exist, got", err)
		return
	}

	commandsFile, err := os.Open(commandsFilePath)

	if err != nil {
		t.Error("Expected to open commands config file, got", err)
		return
	}

	defer commandsFile.Close()

	// Expect commands config file to contain an empty JSON object
	expected := "{}"

	buf := make([]byte, len(expected))

	_, err = commandsFile.Read(buf)

	if err != nil {
		t.Error("Expected to read commands config file, got", err)
		return
	}

	if string(buf) != expected {
		t.Error("Expected commands config file to contain", expected, "got", string(buf))
		return
	}
}

func TestProjectsConfigInit(t *testing.T) {
	s := TestingCore("projects_config_init")

	projectsFilePath := path.Join(s.getConfig().GetConfigDir(), config.ProjectsFileName)

	if _, err := os.Stat(projectsFilePath); os.IsNotExist(err) {
		t.Error("Expected projects to exist, got", err)
		return
	}

	projectsFile, err := os.Open(projectsFilePath)

	if err != nil {
		t.Error("Expected to open projects config file, got", err)
		return
	}

	defer projectsFile.Close()

	// Expect projects config file to contain an empty JSON object
	expected := "{}"

	buf := make([]byte, len(expected))

	_, err = projectsFile.Read(buf)

	if err != nil {
		t.Error("Expected to read projects config file, got", err)
		return
	}

	if string(buf) != expected {
		t.Error("Expected projects config file to contain", expected, "got", string(buf))
		return
	}
}

func TestHostsConfigInit(t *testing.T) {
	s := TestingCore("hosts_config_init")

	hostsFilePath := s.getConfig().GetHostsFile()

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
	s := TestingCore("hosts_backup_dir_init")

	hostsBackupDir := s.getConfig().GetHostsBackupDir()

	if _, err := os.Stat(hostsBackupDir); os.IsNotExist(err) {
		t.Error("Expected hosts backup dir to exist, got", err)
		return
	}
}
