package config

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/iskandervdh/spinup/cli"
)

var beginMarker = fmt.Sprintf("### BEGIN_%s_HOSTS\n", strings.ToUpper(ProgramName))
var endMarker = fmt.Sprintf("\n### END_%s_HOSTS", strings.ToUpper(ProgramName))

func backupHosts() {
	// Create a directory to store the backups if it doesn't exist
	exec.Command("sudo", "mkdir", "-p", "/etc/hosts_bak").Run()

	// Backup the hosts file with a timestamp
	fileName := fmt.Sprintf("/etc/hosts_bak/hosts_%s.bak", time.Now().Format("2006-01-02_15:04:05"))
	exec.Command("sudo", "cp", "/etc/hosts", fileName).Run()
}

func getHostsContent() (string, int, int, error) {
	hosts, err := os.ReadFile("/etc/hosts")

	if err != nil {
		return "", 0, 0, err
	}

	hostsContent := string(hosts)

	beginIndex := strings.Index(hostsContent, beginMarker)
	endIndex := strings.Index(hostsContent, endMarker)

	if beginIndex == -1 || endIndex == -1 || beginIndex >= endIndex {
		return "", beginIndex, endIndex, fmt.Errorf("%s hosts section not found", ProgramName)
	}

	return hostsContent, beginIndex, endIndex, nil
}

func AddHost(domain string) error {
	hostsContent, beginIndex, endIndex, err := getHostsContent()

	if err != nil {
		return err
	}

	customHosts := hostsContent[beginIndex+len(beginMarker) : endIndex]

	// Add domain to hosts file
	customHosts += fmt.Sprintf("\n127.0.0.1\t%s", domain)

	// Save hosts file
	newHostsContent := hostsContent[:beginIndex+len(beginMarker)] + customHosts + hostsContent[endIndex:]

	backupHosts()
	saveNewHosts := exec.Command("sudo", "tee", "/etc/hosts")
	saveNewHosts.Stdin = strings.NewReader(newHostsContent)

	err = saveNewHosts.Run()

	if err != nil {
		return err
	}

	return nil
}

func RemoveHost(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain is empty")
	}

	hostsContent, beginIndex, endIndex, err := getHostsContent()

	if err != nil {
		return err
	}

	customHosts := hostsContent[beginIndex+len(beginMarker) : endIndex]

	// Remove domain to hosts file
	customHosts = strings.Replace(customHosts, fmt.Sprintf("\n127.0.0.1\t%s", domain), "", -1)

	// Save hosts file
	newHostsContent := hostsContent[:beginIndex+len(beginMarker)] + customHosts + hostsContent[endIndex:]

	backupHosts()
	saveNewHosts := exec.Command("sudo", "tee", "/etc/hosts")
	saveNewHosts.Stdin = strings.NewReader(newHostsContent)

	err = saveNewHosts.Run()

	if err != nil {
		return err
	}

	return nil
}

func InitHosts() error {
	hostsContent, beginIndex, endIndex, _ := getHostsContent()

	if beginIndex != -1 && endIndex != -1 {
		cli.WarningPrint("Hosts file already initialized\nSkipping initialization...")
		return nil
	}

	hostsContent = strings.TrimSpace(hostsContent)
	hostsContent += "\n\n"
	hostsContent += beginMarker
	hostsContent += endMarker

	backupHosts()
	saveNewHosts := exec.Command("sudo", "tee", "/etc/hosts")
	saveNewHosts.Stdin = strings.NewReader(hostsContent)

	err := saveNewHosts.Run()

	if err != nil {
		return err
	}

	return nil
}
