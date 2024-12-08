package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/iskandervdh/spinup/cli"
)

func (c *Config) backupHosts() {
	// Create a directory to store the backups if it doesn't exist
	exec.Command("sudo", "mkdir", "-p", hostsBackupDir).Run()

	// Backup the hosts file with a timestamp
	fileName := fmt.Sprintf("%s/hosts_%s.bak", hostsBackupDir, time.Now().Format("2006-01-02_15:04:05"))
	exec.Command("sudo", "cp", hostsFile, fileName).Run()
}

func (c *Config) getHostsContent() (string, int, int, error) {
	hosts, err := os.ReadFile(hostsFile)

	if err != nil {
		return "", 0, 0, err
	}

	hostsContent := string(hosts)

	beginIndex := strings.Index(hostsContent, c.hostsBeginMarker)
	endIndex := strings.Index(hostsContent, c.hostsEndMarker)

	if beginIndex == -1 || endIndex == -1 || beginIndex >= endIndex {
		return "", beginIndex, endIndex, fmt.Errorf("%s hosts section not found", ProgramName)
	}

	return hostsContent, beginIndex, endIndex, nil
}

func (c *Config) AddHost(domain string) error {
	hostsContent, beginIndex, endIndex, err := c.getHostsContent()

	if err != nil {
		return err
	}

	customHosts := hostsContent[beginIndex+len(c.hostsBeginMarker) : endIndex]

	// Add domain to hosts file
	customHosts += fmt.Sprintf("\n127.0.0.1\t%s", domain)

	// Save hosts file
	newHostsContent := hostsContent[:beginIndex+len(c.hostsBeginMarker)] + customHosts + hostsContent[endIndex:]

	c.backupHosts()
	saveNewHosts := exec.Command("sudo", "tee", hostsFile)
	saveNewHosts.Stdin = strings.NewReader(newHostsContent)

	err = saveNewHosts.Run()

	if err != nil {
		return err
	}

	return nil
}

func (c *Config) RemoveHost(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain is empty")
	}

	hostsContent, beginIndex, endIndex, err := c.getHostsContent()

	if err != nil {
		return err
	}

	customHosts := hostsContent[beginIndex+len(c.hostsBeginMarker) : endIndex]

	// Remove domain to hosts file
	customHosts = strings.Replace(customHosts, fmt.Sprintf("\n127.0.0.1\t%s", domain), "", -1)

	// Save hosts file
	newHostsContent := hostsContent[:beginIndex+len(c.hostsBeginMarker)] + customHosts + hostsContent[endIndex:]

	c.backupHosts()
	saveNewHosts := exec.Command("sudo", "tee", hostsFile)
	saveNewHosts.Stdin = strings.NewReader(newHostsContent)

	err = saveNewHosts.Run()

	if err != nil {
		return err
	}

	return nil
}

func (c *Config) InitHosts() error {
	hostsContent, beginIndex, endIndex, _ := c.getHostsContent()

	if beginIndex != -1 && endIndex != -1 {
		cli.WarningPrint("Hosts file already initialized\nSkipping initialization...")
		return nil
	}

	hostsContent = strings.TrimSpace(hostsContent)
	hostsContent += "\n\n"
	hostsContent += c.hostsBeginMarker
	hostsContent += c.hostsEndMarker

	c.backupHosts()
	saveNewHosts := exec.Command("sudo", "tee", hostsFile)
	saveNewHosts.Stdin = strings.NewReader(hostsContent)

	err := saveNewHosts.Run()

	if err != nil {
		return err
	}

	return nil
}
