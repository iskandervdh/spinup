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
