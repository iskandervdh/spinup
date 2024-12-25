package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func (c *Config) restartNginx() error {
	return exec.Command("sudo", "systemctl", "restart", "nginx").Run()
}

func (c *Config) AddNginxConfig(name string, domain string, port int) error {
	config := fmt.Sprintf(`server {
	listen 80;

	server_name %s;

	location / {
		proxy_pass http://127.0.0.1:%d/;
		proxy_set_header Host $host;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header X-Forwarded-Proto $scheme;
	}
}
`, domain, port)

	nginxConfigFilePath := fmt.Sprintf("%s/%s.conf", c.nginxConfigDir, name)

	if _, err := os.Stat(nginxConfigFilePath); err == nil {
		return fmt.Errorf("config file %s already exists", nginxConfigFilePath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if config file exists: %v", err)
	}

	err := c.withSudo("touch", nginxConfigFilePath).Run()

	if err != nil {
		return err
	}

	createCommand := c.withSudo("tee", nginxConfigFilePath)
	createCommand.Stdin = strings.NewReader(config)
	err = createCommand.Run()

	if err != nil {
		return err
	}

	if !c.IsTesting() {
		c.restartNginx()
	}

	return nil
}

func (c *Config) RemoveNginxConfig(name string) error {
	nginxConfigFilePath := fmt.Sprintf("%s/%s.conf", c.nginxConfigDir, name)
	err := c.withSudo("rm", nginxConfigFilePath).Run()

	if err != nil {
		return err
	}

	if !c.IsTesting() {
		c.restartNginx()
	}

	return nil
}

func (c *Config) UpdateNginxConfig(name string, domain string, port int) error {
	err := c.RemoveNginxConfig(name)

	if err != nil {
		return err
	}

	return c.AddNginxConfig(name, domain, port)
}

func (c *Config) RenameNginxConfig(oldName string, newName string) error {
	oldNginxConfigFilePath := fmt.Sprintf("%s/%s.conf", c.nginxConfigDir, oldName)
	newNginxConfigFilePath := fmt.Sprintf("%s/%s.conf", c.nginxConfigDir, newName)

	err := c.withSudo("mv", oldNginxConfigFilePath, newNginxConfigFilePath).Run()

	if err != nil {
		return err
	}

	if !c.IsTesting() {
		c.restartNginx()
	}

	return nil
}

func (c *Config) InitNginx() error {
	if _, err := os.Stat(c.nginxConfigDir); os.IsNotExist(err) {
		err := os.MkdirAll(c.nginxConfigDir, 0755)

		if err != nil {
			return err
		}
	}

	return nil
}
