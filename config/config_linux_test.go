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
