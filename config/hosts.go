package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/iskandervdh/spinup/cli"
)

var HostsBeginMarker = fmt.Sprintf("### BEGIN_%s_HOSTS", strings.ToUpper(ProgramName))
var HostsEndMarker = fmt.Sprintf("\n### END_%s_HOSTS", strings.ToUpper(ProgramName))

func (c *Config) backupHosts() error {
	// Create a directory to store the backups if it doesn't exist
	err := c.withSudo("mkdir", "-p", c.hostsBackupDir).Run()

	if err != nil {
		return err
	}

	// Backup the hosts file with a timestamp
	fileName := fmt.Sprintf("%s/hosts_%s.bak", c.hostsBackupDir, time.Now().Format("2006-01-02_15:04:05"))
	err = c.withSudo("cp", c.hostsFile, fileName).Run()

	if err != nil {
		return err
	}

	return nil
}

func (c *Config) getHostsContent() (string, int, int, error) {
	hosts, err := os.ReadFile(c.hostsFile)

	if err != nil {
		return "", 0, 0, err
	}

	hostsContent := string(hosts)

	beginIndex := strings.Index(hostsContent, HostsBeginMarker)
	endIndex := strings.Index(hostsContent, HostsEndMarker)

	if beginIndex == -1 || endIndex == -1 || beginIndex >= endIndex {
		return "", beginIndex, endIndex, fmt.Errorf("%s hosts section not found", ProgramName)
	}

	return hostsContent, beginIndex, endIndex, nil
}

func (c *Config) AddHost(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain is empty")
	}

	hostsContent, beginIndex, endIndex, err := c.getHostsContent()

	if err != nil {
		return err
	}

	customHosts := hostsContent[beginIndex+len(HostsBeginMarker) : endIndex]

	// Check if host already exists
	if strings.Contains(customHosts, fmt.Sprintf("\n127.0.0.1\t%s", domain)) {
		return fmt.Errorf("domain %s already exists", domain)
	}

	// Add domain to hosts file
	customHosts += fmt.Sprintf("\n127.0.0.1\t%s", domain)

	// Save hosts file
	newHostsContent := hostsContent[:beginIndex+len(HostsBeginMarker)] + customHosts + hostsContent[endIndex:]

	err = c.backupHosts()

	if err != nil {
		return err
	}

	saveNewHosts := c.withSudo("tee", c.hostsFile)
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

	customHosts := hostsContent[beginIndex+len(HostsBeginMarker) : endIndex]

	// Remove domain to hosts file
	customHosts = strings.Replace(customHosts, fmt.Sprintf("\n127.0.0.1\t%s", domain), "", -1)

	// Save hosts file
	newHostsContent := hostsContent[:beginIndex+len(HostsBeginMarker)] + customHosts + hostsContent[endIndex:]

	c.backupHosts()
	saveNewHosts := c.withSudo("tee", c.hostsFile)
	saveNewHosts.Stdin = strings.NewReader(newHostsContent)

	err = saveNewHosts.Run()

	if err != nil {
		return err
	}

	return nil
}

func (c *Config) InitHosts(cli *cli.CLI) error {
	// Check if hosts file exists
	fileInfo, err := os.Stat(c.hostsFile)

	// Create hosts file if it doesn't exist
	if err != nil {
		if os.IsNotExist(err) {
			_, err := os.Create(c.hostsFile)

			if err != nil {
				return err
			}

			// Check if hosts file was created
			fileInfo, err = os.Stat(c.hostsFile)

			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("hosts file is a directory")
	}

	hostsContent, beginIndex, endIndex, _ := c.getHostsContent()

	if beginIndex != -1 && endIndex != -1 {
		cli.WarningPrint("Hosts file already initialized\nSkipping initialization...")
		return nil
	}

	hostsContent = strings.TrimSpace(hostsContent)
	hostsContent += "\n\n"
	hostsContent += HostsBeginMarker
	hostsContent += HostsEndMarker

	err = c.backupHosts()

	if err != nil {
		return err
	}

	saveNewHosts := c.withSudo("tee", c.hostsFile)
	saveNewHosts.Stdin = strings.NewReader(hostsContent)

	err = saveNewHosts.Run()

	if err != nil {
		return err
	}

	return nil
}
