//go:build windows

package config

import (
	"io"
	"os"
	"path/filepath"
)

func getHostsFilePath() string {
	home, err := os.UserHomeDir()

	if err != nil {
		return ""
	}

	return filepath.Join(home, "Documents", "hosts")
}

func getHostsBackupDir() string {
	home, err := os.UserHomeDir()

	if err != nil {
		return ""
	}

	return filepath.Join(home, "Documents", "hosts_bak")
}

var nginxConfigDir = "C:\\nginx\\conf\\conf.d"
var hostsFile = getHostsFilePath()
var hostsBackupDir = getHostsBackupDir()

func (c *Config) createHostsBackupDir() error {
	return os.MkdirAll(c.hostsBackupDir, os.ModeDir)
}

func (c *Config) copyFile(from string, to string) error {
	r, err := os.Open(from)

	if err != nil {
		return err
	}

	defer r.Close()

	w, err := os.Create(to)

	if err != nil {
		return err
	}

	defer func() {
		// Only return the error of closing the file if there was no error before
		if c := w.Close(); err == nil {
			err = c
		}
	}()

	_, err = io.Copy(w, r)

	return err
}

func (c *Config) writeToFile(filePath string, contents string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return err
	}

	_, err = f.Write([]byte(contents))

	if err != nil {
		return err
	}

	return f.Close()
}
