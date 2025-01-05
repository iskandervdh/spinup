package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/iskandervdh/spinup/common"
)

// Markers that indicate the beginning and end of the custom hosts section.
var HostsBeginMarker = fmt.Sprintf("### BEGIN_%s_HOSTS", strings.ToUpper(common.ProgramName))
var HostsEndMarker = fmt.Sprintf("\n### END_%s_HOSTS", strings.ToUpper(common.ProgramName))

// Backup the hosts file with a timestamp so in case of an error the original file can be restored.
func (c *Config) backupHosts() error {
	// Create a directory to store the backups if it doesn't exist
	err := c.createHostsBackupDir()

	if err != nil {
		return err
	}

	// Backup the hosts file with a timestamp
	fileName := filepath.Join(c.hostsBackupDir, fmt.Sprintf("hosts_%s.bak", time.Now().Format("2006-01-02_15-04-05")))
	err = c.copyFile(c.hostsFile, fileName)

	if err != nil {
		return err
	}

	return nil
}

// Get the content of the hosts file and the begin and end of the custom hosts section.
func (c *Config) getHostsContent() (string, int, int, error) {
	hosts, err := os.ReadFile(c.hostsFile)

	if err != nil {
		return "", 0, 0, err
	}

	hostsContent := string(hosts)

	beginIndex := strings.Index(hostsContent, HostsBeginMarker)
	endIndex := strings.Index(hostsContent, HostsEndMarker)

	if beginIndex == -1 || endIndex == -1 || beginIndex >= endIndex {
		return "", beginIndex, endIndex, fmt.Errorf("%s hosts section not found", common.ProgramName)
	}

	return hostsContent, beginIndex, endIndex, nil
}

// Add a domain to the custom hosts section.
func (c *Config) AddDomain(domain string) error {
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

	return c.writeToFile(c.hostsFile, newHostsContent)
}

// Remove a domain from the custom hosts section.
func (c *Config) RemoveDomain(domain string) error {
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

	return c.writeToFile(c.hostsFile, newHostsContent)
}

// Update a domain in the custom hosts section.
func (c *Config) UpdateHost(oldDomain string, newDomain string) error {
	if oldDomain == "" {
		return fmt.Errorf("old domain is empty")
	}

	if newDomain == "" {
		return fmt.Errorf("new domain is empty")
	}

	hostsContent, beginIndex, endIndex, err := c.getHostsContent()

	if err != nil {
		return err
	}

	customHosts := hostsContent[beginIndex+len(HostsBeginMarker) : endIndex]

	// Update domain in hosts file
	customHosts = strings.Replace(customHosts, fmt.Sprintf("\n127.0.0.1\t%s", oldDomain), fmt.Sprintf("\n127.0.0.1\t%s", newDomain), -1)

	// Save hosts file
	newHostsContent := hostsContent[:beginIndex+len(HostsBeginMarker)] + customHosts + hostsContent[endIndex:]

	c.backupHosts()

	return c.writeToFile(c.hostsFile, newHostsContent)
}

// Initialize the custom hosts section in the hosts file
// so domains can be added without overwriting the other contents of the file.
func (c *Config) InitHosts() error {
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
		fmt.Println("Hosts file already initialized\nSkipping initialization...")
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

	err = c.writeToFile(c.hostsFile, hostsContent)

	if err != nil {
		return err
	}

	return nil
}
