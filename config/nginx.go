package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/iskandervdh/spinup/common"
)

// Regex to match the server_name directive in a Nginx config file.
var serverNameRegex = regexp.MustCompile(`server_name\s+(.*);`)

// Reload the Nginx service.
func (c *Config) reloadNginx() error {
	return exec.Command("sudo", "systemctl", "reload", "nginx").Run()
}

// Add a new Nginx configuration file with the given name and port.
func (c *Config) AddNginxConfig(name string, port int64) error {
	config := fmt.Sprintf(`server {
	listen 80;

	server_name %s.test;

	location / {
		proxy_pass http://127.0.0.1:%d/;
		proxy_set_header Host $host;
		proxy_set_header X-Real-IP $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
		proxy_set_header X-Forwarded-Proto $scheme;
	}
}
`, name, port)

	nginxConfigFilePath := fmt.Sprintf("%s/%s.conf", c.nginxConfigDir, name)

	if _, err := os.Stat(nginxConfigFilePath); err == nil {
		return fmt.Errorf("config file %s already exists", nginxConfigFilePath)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check if config file exists: %v", err)
	}

	err := c.writeToFile(nginxConfigFilePath, config)

	if err != nil {
		return err
	}

	if !c.IsTesting() {
		c.reloadNginx()
	}

	return nil
}

// Remove a Nginx configuration file with the given name.
func (c *Config) RemoveNginxConfig(name string) error {
	nginxConfigFilePath := fmt.Sprintf("%s/%s.conf", c.nginxConfigDir, name)
	err := os.Remove(nginxConfigFilePath)

	if err != nil {
		return err
	}

	if !c.IsTesting() {
		c.reloadNginx()
	}

	return nil
}

// Update a Nginx configuration file with the given name and port.
func (c *Config) UpdateNginxConfig(name string, port int64) error {
	err := c.RemoveNginxConfig(name)

	if err != nil {
		return err
	}

	return c.AddNginxConfig(name, port)
}

// Rename a Nginx configuration file with the given old and new name.
func (c *Config) RenameNginxConfig(oldName string, newName string) error {
	oldNginxConfigFilePath := fmt.Sprintf("%s/%s.conf", c.nginxConfigDir, oldName)
	newNginxConfigFilePath := fmt.Sprintf("%s/%s.conf", c.nginxConfigDir, newName)

	err := os.Rename(oldNginxConfigFilePath, newNginxConfigFilePath)

	if err != nil {
		return err
	}

	if !c.IsTesting() {
		c.reloadNginx()
	}

	return nil
}

// Add a domain alias to a Nginx configuration file.
func (c *Config) NginxAddDomainAlias(name string, domainAlias string) error {
	nginxConfigFilePath := fmt.Sprintf("%s/%s.conf", c.nginxConfigDir, name)
	content, err := os.ReadFile(nginxConfigFilePath)

	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	serverName := ""

	for _, line := range lines {
		if serverNameRegex.MatchString(line) {
			serverName = serverNameRegex.FindStringSubmatch(line)[1]
		}
	}

	if serverName == "" {
		return fmt.Errorf("server_name not found in config file")
	}

	newServerName := fmt.Sprintf("server_name %s %s;", serverName, domainAlias)
	updatedConfig := strings.ReplaceAll(string(content), fmt.Sprintf("server_name %s;", serverName), newServerName)

	err = c.writeToFile(nginxConfigFilePath, updatedConfig)

	if err != nil {
		return err
	}

	if !c.IsTesting() {
		c.reloadNginx()
	}

	return nil
}

// Remove a domain alias from a Nginx configuration file.
func (c *Config) NginxRemoveDomainAlias(name string, domainAlias string) error {
	nginxConfigFilePath := fmt.Sprintf("%s/%s.conf", c.nginxConfigDir, name)
	content, err := os.ReadFile(nginxConfigFilePath)

	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	serverName := ""

	for _, line := range lines {
		if serverNameRegex.MatchString(line) {
			serverName = serverNameRegex.FindStringSubmatch(line)[1]
		}
	}

	if serverName == "" {
		return fmt.Errorf("server_name not found in config file")
	}

	updatedServerName := strings.Trim(strings.ReplaceAll(serverName, domainAlias, ""), " ")
	newServerName := fmt.Sprintf("server_name %s;", updatedServerName)
	updatedConfig := strings.ReplaceAll(string(content), fmt.Sprintf("server_name %s;", serverName), newServerName)

	err = c.writeToFile(nginxConfigFilePath, updatedConfig)

	if err != nil {
		return err
	}

	if !c.IsTesting() {
		c.reloadNginx()
	}

	return nil
}

// Initialize the Nginx configuration directory.
func (c *Config) InitNginx() error {
	if _, err := os.Stat(c.nginxConfigDir); os.IsNotExist(err) {
		err := os.MkdirAll(c.nginxConfigDir, 0755)

		if err != nil {
			return err
		}
	}

	if common.IsWindows() && !c.IsTesting() {
		fmt.Printf(
			"\n!!! Please add the following include directive to http section of your nginx.conf file located at %s like this:\n",
			filepath.Join(c.nginxConfigDir, "..", "nginx.conf"),
		)
		fmt.Printf("\nhttp {\n\t...\n\n\t%s\n}\n", "include \"C:/nginx/conf/conf.d/*.conf\";")
	} else if !c.IsTesting() {
		fmt.Printf(
			"\n!!! Please add the following include directive to http section of your nginx.conf file located at %s like this:\n",
			filepath.Join(c.nginxConfigDir, "..", "nginx.conf"),
		)
		fmt.Printf("\nhttp {\n\t...\n\n\t%s\n}\n", "include /etc/spinup/config/nginx/**/*.conf;")
	}

	return nil
}
