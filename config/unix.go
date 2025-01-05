//go:build linux || darwin

package config

import "strings"

var nginxConfigDir = "/etc/nginx/conf.d"
var hostsFile = "/etc/hosts"
var hostsBackupDir = "/etc/hosts_bak"

func (c *Config) createHostsBackupDir() error {
	return c.withSudo("mkdir", "-p", c.hostsBackupDir).Run()
}

func (c *Config) copyFile(from string, to string) error {
	return c.withSudo("cp", from, to).Run()
}

func (c *Config) writeToFile(filePath string, content string) error {
	saveNewHosts := c.withSudo("tee", filePath)
	saveNewHosts.Stdin = strings.NewReader(content)

	return saveNewHosts.Run()
}
