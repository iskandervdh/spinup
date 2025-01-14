//go:build windows

package config

import (
	"os"
)

var nginxConfigDir = "C:\\nginx\\conf\\conf.d"

func (c *Config) writeToFile(filePath string, contents string) error {
	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		return err
	}

	_, err = f.Write([]byte(contents))

	if err != nil {
		return err
	}

	return f.Close()
}

func (c *Config) moveFile(oldFilePath string, newFilePath string) error {
	return os.Rename(oldFilePath, newFilePath)
}

func (c *Config) removeFile(filePath string) error {
	return os.Remove(filePath)
}
