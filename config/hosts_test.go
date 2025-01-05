package config

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/iskandervdh/spinup/common"
)

func TestInitHosts(t *testing.T) {
	c := TestingConfig("init_hosts")

	err := c.InitHosts()

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

	err = c.InitHosts()

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestInitHostsAlreadyInitialized(t *testing.T) {
	c := TestingConfig("init_hosts_already_initialized")

	r, w, err := os.Pipe()

	stdout := os.Stdout
	stderr := os.Stderr

	os.Stdout = w
	os.Stderr = w

	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
	}()

	if err != nil {
		t.Fatal(err)
	}

	err = c.InitHosts()

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	err = c.InitHosts()

	if err != nil {
		t.Error("Expected error, got nil")
	}

	err = c.InitHosts()

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

	err := c.InitHosts()

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

	err = c.InitHosts()

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestBackupHosts(t *testing.T) {
	c := TestingConfig("backup_hosts")
	c.InitHosts()

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

func TestAddDomain(t *testing.T) {
	c := TestingConfig("add_domain")
	c.InitHosts()

	err := c.AddDomain("test.local")

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	hostsContent, _, _, _ := c.getHostsContent()

	expected := "\n\n" + HostsBeginMarker + "\n127.0.0.1\ttest.local" + HostsEndMarker

	if common.IsWindows() {
		hostsContent = strings.ReplaceAll(hostsContent, "\u0000", "")
	}

	if hostsContent != expected {
		t.Errorf("Expected %s, got %s", expected, hostsContent)
	}
}

func TestAddDomainDuplicate(t *testing.T) {
	c := TestingConfig("add_host_duplicate")
	c.InitHosts()

	err := c.AddDomain("test.local")

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	err = c.AddDomain("test.local")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAddDomainFileError(t *testing.T) {
	c := TestingConfig("add_host_file_error")

	err := c.AddDomain("test.local")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAddDomainPermissionError(t *testing.T) {
	c := TestingConfig("add_host_permission_error")
	c.InitHosts()

	os.Chmod(c.GetHostsFile(), 0444)

	err := c.AddDomain("test.local")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestAddDomainEmpty(t *testing.T) {
	c := TestingConfig("add_host_empty")
	c.InitHosts()

	err := c.AddDomain("")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestRemoveDomain(t *testing.T) {
	c := TestingConfig("remove_host")
	c.InitHosts()
	c.AddDomain("test.local")

	err := c.RemoveDomain("test.local")

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	hostsContent, _, _, _ := c.getHostsContent()

	expected := "\n\n" + HostsBeginMarker + HostsEndMarker

	if common.IsWindows() {
		hostsContent = strings.ReplaceAll(hostsContent, "\u0000", "")
	}

	if hostsContent != expected {
		t.Errorf("Expected %s, got %s", expected, hostsContent)
	}
}

func TestRemoveDomainEmpty(t *testing.T) {
	c := TestingConfig("remove_host_empty")
	c.InitHosts()

	err := c.RemoveDomain("")

	if err == nil {
		t.Error("Expected error, got nil")
	}
}
