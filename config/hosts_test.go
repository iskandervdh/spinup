package config

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/iskandervdh/spinup/cli"
)

func TestInitHosts(t *testing.T) {
	c := TestingConfig("init_hosts")

	err := c.InitHosts(cli.New())

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestInitHostsFileError(t *testing.T) {
	c := TestingConfig("init_hosts_file_error")

	testingConfigDir := TestingConfigDir("init_hosts_file_error")
	err := os.RemoveAll(testingConfigDir)

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	err = c.InitHosts(cli.New())

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestInitHostsAlreadyInitialized(t *testing.T) {
	c := TestingConfig("init_hosts_already_initialized")

	r, w, err := os.Pipe()

	if err != nil {
		t.Fatal(err)
	}

	cli := cli.New(cli.WithOut(w))
	err = c.InitHosts(cli)

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	err = c.InitHosts(cli)

	if err != nil {
		t.Error("Expected error, got nil")
	}

	err = c.InitHosts(cli)

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	w.Close()
	out, err := io.ReadAll(r)

	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(string(out), "already initialized") {
		t.Errorf("Expected warning message, got none %s", string(out))
	}
}

func TestInitHostsPermissionError(t *testing.T) {
	c := TestingConfig("init_hosts_permission_error")

	os.Create(c.GetHostsFile())
	os.Chmod(c.GetHostsFile(), 0444)

	err := c.InitHosts(cli.New())

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestInitHostsStatError(t *testing.T) {
	c := TestingConfig("init_hosts_stat_error")

	err := os.MkdirAll(c.GetHostsFile(), 0755)

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	err = c.InitHosts(cli.New())

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestInitHostsBackupError(t *testing.T) {
	c := TestingConfig("init_hosts_backup_error")

	os.Mkdir(c.GetHostsBackupDir(), 0444)

	err := c.InitHosts(cli.New())

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestBackupHosts(t *testing.T) {
	c := TestingConfig("backup_hosts")
	c.InitHosts(cli.New())

	err := c.backupHosts()

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestBackupHostsDirError(t *testing.T) {
	c := TestingConfig("backup_hosts_dir_error")
	os.Chmod(c.GetConfigDir(), 0444)

	err := c.backupHosts()

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestBackupHostsFileError(t *testing.T) {
	c := TestingConfig("backup_hosts_file_error")

	err := c.backupHosts()

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetHostsContentFileError(t *testing.T) {
	c := TestingConfig("get_hosts_content_file_error")

	content, beginIndex, endIndex, err := c.getHostsContent()

	if content != "" {
		t.Errorf("Expected empty content, got %s", content)
	}

	if beginIndex != 0 {
		t.Errorf("Expected beginIndex to be 0, got %d", beginIndex)
	}

	if endIndex != 0 {
		t.Errorf("Expected endIndex to be 0, got %d", endIndex)
	}

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAddHost(t *testing.T) {
	c := TestingConfig("add_host")
	c.InitHosts(cli.New())

	err := c.AddHost("test.local")

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	content, _, _, _ := c.getHostsContent()

	expected := "\n\n" + HostsBeginMarker + "\n127.0.0.1\ttest.local" + HostsEndMarker

	if content != expected {
		t.Errorf("Expected %s, got %s", expected, content)
	}
}

func TestAddHostDuplicate(t *testing.T) {
	c := TestingConfig("add_host_duplicate")
	c.InitHosts(cli.New())

	err := c.AddHost("test.local")

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	err = c.AddHost("test.local")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAddHostFileError(t *testing.T) {
	c := TestingConfig("add_host_file_error")

	err := c.AddHost("test.local")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAddHostPermissionError(t *testing.T) {
	c := TestingConfig("add_host_permission_error")
	c.InitHosts(cli.New())

	os.Chmod(c.GetHostsFile(), 0444)

	err := c.AddHost("test.local")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAddHostEmpty(t *testing.T) {
	c := TestingConfig("add_host_empty")
	c.InitHosts(cli.New())

	err := c.AddHost("")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestRemoveHost(t *testing.T) {
	c := TestingConfig("remove_host")
	c.InitHosts(cli.New())
	c.AddHost("test.local")

	err := c.RemoveHost("test.local")

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	content, _, _, _ := c.getHostsContent()

	expected := "\n\n" + HostsBeginMarker + HostsEndMarker

	if content != expected {
		t.Errorf("Expected %s, got %s", expected, content)
	}
}

func TestRemoveHostEmpty(t *testing.T) {
	c := TestingConfig("remove_host_empty")
	c.InitHosts(cli.New())

	err := c.RemoveHost("")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}
