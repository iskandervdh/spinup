package spinup

import (
	"os"
	"path"
	"testing"

	"github.com/iskandervdh/spinup/cli"
	"github.com/iskandervdh/spinup/config"
)

func TestingSpinupConfigDir(testName string) string {
	return path.Join(os.TempDir(), config.ProgramName, testName)
}

func TestingSpinup(testName string, _cli *cli.CLI) *Spinup {
	// Remove old tmp config dir
	testingConfigDir := TestingSpinupConfigDir(testName)
	err := os.RemoveAll(testingConfigDir)

	if err != nil {
		panic(err)
	}

	if _cli == nil {
		_cli = cli.New()
	}

	c := config.NewTesting(testingConfigDir)
	s := New(WithConfig(c), WithCLI(_cli))
	s.init()

	return s
}

func TestNew(t *testing.T) {
	testName := "new"
	s := TestingSpinup(testName, nil)

	if s == nil {
		t.Error("Expected Spinup instance, got nil")
		return
	}

	if s.getConfig() == nil {
		t.Error("Expected Config instance, got nil")
		return
	}

	if s.getConfig().GetConfigDir() != TestingSpinupConfigDir(testName) {
		t.Error(
			"Expected ConfigDir to be",
			TestingSpinupConfigDir(testName),
			"got",
			s.getConfig().GetConfigDir(),
		)
	}
}

func TestInit(t *testing.T) {
	s := TestingSpinup("init", nil)

	configDir := s.getConfig().GetConfigDir()

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("Expected ConfigDir to exist, got", err)
	}
}

func TestCommandsConfigInit(t *testing.T) {
	s := TestingSpinup("commands_config_init", nil)

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
	s := TestingSpinup("projects_config_init", nil)

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
	s := TestingSpinup("hosts_config_init", nil)

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
	s := TestingSpinup("hosts_backup_dir_init", nil)

	hostsBackupDir := s.getConfig().GetHostsBackupDir()

	if _, err := os.Stat(hostsBackupDir); os.IsNotExist(err) {
		t.Error("Expected hosts backup dir to exist, got", err)
		return
	}
}

func TestSpinupHandleUnknownSubcommand(t *testing.T) {
	s := TestingSpinup("handle", nil)

	// Test handle without any arguments
	os.Args = []string{"spinup", "handle"}
	s.Handle()
}

func TestSpinupHandleNoArgs(t *testing.T) {
	s := TestingSpinup("handle_no_args", nil)

	os.Args = []string{"spinup"}
	s.Handle()
}

func TestSpinupHandleInit(t *testing.T) {
	s := TestingSpinup("handle_no_args", nil)

	os.Args = []string{"spinup", "init"}
	s.Handle()
}

func TestSpinupHandleVersion(t *testing.T) {
	s := TestingSpinup("handle_no_args", nil)

	os.Args = []string{"spinup", "-v"}
	s.Handle()
}

func TestSpinupHandleCommand(*testing.T) {
	s := TestingSpinup("handle_no_args", nil)

	os.Args = []string{"spinup", "c"}
	s.Handle()

	os.Args = []string{"spinup", "c", "ls"}
	s.Handle()

	os.Args = []string{"spinup", "c", "add", "test"}
	s.Handle()

	os.Args = []string{"spinup", "c", "add", "test", "echo test"}
	s.Handle()

	os.Args = []string{"spinup", "c", "rm", "test"}
	s.Handle()

	os.Args = []string{"spinup", "c", "test"}
	s.Handle()
}

func TestSpinupHandleProject(t *testing.T) {
	s := TestingSpinup("handle_no_args", nil)

	os.Args = []string{"spinup", "p"}
	s.Handle()

	os.Args = []string{"spinup", "p", "ls"}
	s.Handle()

	os.Args = []string{"spinup", "p", "add", "test"}
	s.Handle()

	os.Args = []string{"spinup", "p", "add", "test", "echo test"}
	s.Handle()

	os.Args = []string{"spinup", "p", "rm", "test"}
	s.Handle()

	os.Args = []string{"spinup", "p", "test"}
	s.Handle()
}

func TestSpinupHandleVariable(t *testing.T) {
	s := TestingSpinup("handle_no_args", nil)

	os.Args = []string{"spinup", "v"}
	s.Handle()

	os.Args = []string{"spinup", "v", "ls"}
	s.Handle()

	os.Args = []string{"spinup", "v", "ls", "test"}
	s.Handle()

	os.Args = []string{"spinup", "v", "add", "test"}
	s.Handle()

	os.Args = []string{"spinup", "v", "add", "test", "echo test"}
	s.Handle()

	os.Args = []string{"spinup", "v", "rm", "test"}
	s.Handle()

	os.Args = []string{"spinup", "v", "test"}
	s.Handle()
}

func TestSpinupHandle(t *testing.T) {
	s := TestingSpinup("handle_no_args", nil)

	os.Args = []string{"spinup", "run", "test"}
	s.Handle()
}
