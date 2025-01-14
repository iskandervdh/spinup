//go:build linux || darwin

package config

import "strings"

var nginxConfigDir = "/etc/nginx/conf.d"

func (c *Config) writeToFile(filePath string, content string) error {
	saveNewHosts := c.withSudo("tee", filePath)
	saveNewHosts.Stdin = strings.NewReader(content)

	return saveNewHosts.Run()
}

func (c *Config) moveFile(oldFilePath string, newFilePath string) error {
	return c.withSudo("mv", oldFilePath, newFilePath).Run()
}

func (c *Config) removeFile(filePath string) error {
	return c.withSudo("rm", filePath).Run()
}
