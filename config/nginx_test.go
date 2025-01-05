package config

import (
	"testing"
)

func TestInitNginx(t *testing.T) {
	c := TestingConfig("init_nginx")

	err := c.InitNginx()

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestAddNginxConfig(t *testing.T) {
	c := TestingConfig("add_nginx_config")

	err := c.InitNginx()

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	err = c.AddNginxConfig("test", "test.local", 8080)

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestAddNginxConfigExists(t *testing.T) {
	c := TestingConfig("add_nginx_config_exists")

	err := c.InitNginx()

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	err = c.AddNginxConfig("test", "test.local", 8080)

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	err = c.AddNginxConfig("test", "test.local", 8080)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestRemoveNginxConfig(t *testing.T) {
	c := TestingConfig("remove_nginx_config")

	err := c.InitNginx()

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	err = c.AddNginxConfig("test", "test.local", 8080)

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	err = c.RemoveNginxConfig("test")

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
}

func TestRemoveNginxConfigDoesNotExist(t *testing.T) {
	c := TestingConfig("remove_nginx_config_does_not_exist")

	err := c.InitNginx()

	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	err = c.RemoveNginxConfig("test")

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
