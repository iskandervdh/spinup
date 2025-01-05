package config

import (
	"os"
	"testing"
)

func TestNewError(t *testing.T) {
	home, _ := os.LookupEnv("HOME")
	os.Unsetenv("HOME")

	_, err := New()

	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	os.Setenv("HOME", home)
}

func TestWithSudo(t *testing.T) {
	c, err := New()

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	cmd := c.withSudo("ls")

	if cmd.Args[0] != "sudo" {
		t.Errorf("Expected sudo, got %s", cmd.Args[0])
	}
}

func TestInitHostsBackupError(t *testing.T) {
	c := TestingConfig("init_hosts_backup_error")

	os.Mkdir(c.GetHostsBackupDir(), 0444)

	err := c.InitHosts()

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
