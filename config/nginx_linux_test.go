package config

import (
	"os"
	"testing"
)

func TestAddNginxConfigPermissionError(t *testing.T) {
	c := TestingConfig("add_nginx_config_permission_error")

	err := c.InitNginx()

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	os.Chmod(c.nginxConfigDir, 0444)

	err = c.AddNginxConfig("test", 8080)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
