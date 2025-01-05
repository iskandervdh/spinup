package core

import (
	"os"
	"strings"
	"testing"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/config"
)

func TestAddDomainAlias(t *testing.T) {
	c := TestingCore("add_domain_alias")

	// Fetch the commands and projects from their config files
	c.FetchCommands()
	c.FetchProjects()

	c.AddProject("test", "test.local", 1234, []string{})

	// "Refetch" the projects from the config file
	c.FetchProjects()

	c.AddDomainAlias("test", "test.test")

	// "Refetch" the projects from the config file
	c.FetchProjects()

	exists, project := c.ProjectExists("test")

	if !exists {
		t.Error("Expected project to exist")
		return
	}

	if project.DomainAliases == nil {
		t.Error("Expected domain aliases to be initialized")
		return
	}

	if len(project.DomainAliases) != 1 {
		t.Error("Expected 1 domain alias, got", len(project.DomainAliases))
		return
	}

	if project.DomainAliases[0].Value != "test.test" {
		t.Error("Expected domain alias to be 'test.test', got", project.DomainAliases[0])
	}

	// Check if the domain alias is added to the hosts file
	hostsFilePath := c.GetConfig().GetHostsFile()
	hostsFile, err := os.Open(hostsFilePath)

	if err != nil {
		t.Error("Error to open hosts file, got:", err)
		return
	}

	defer hostsFile.Close()

	fileInfo, err := os.Stat(hostsFilePath)

	if err != nil {
		t.Error("Expected to stat hosts file, got", err)
		return
	}

	buf := make([]byte, fileInfo.Size())

	_, err = hostsFile.Read(buf)

	if err != nil {
		t.Error("Expected to read hosts file, got", err)
		return
	}

	hostsContent := string(buf)
	expected := "\n\n" + config.HostsBeginMarker + "\n127.0.0.1\ttest.local\n127.0.0.1\ttest.test" + config.HostsEndMarker

	if common.IsWindows() {
		hostsContent = strings.ReplaceAll(hostsContent, "\u0000", "")
	}

	if hostsContent != expected {
		t.Error("Expected hosts file to contain", expected, "got", hostsContent)
		return
	}
}

func TestRemoveDomainAlias(t *testing.T) {
	c := TestingCore("remove_domain_alias")

	// Fetch the commands and projects from their config files
	c.FetchCommands()
	c.FetchProjects()

	c.AddProject("test", "test.local", 1234, []string{})

	err := c.FetchProjects()

	if err != nil {
		t.Error("Expected to fetch projects, got", err)
		return
	}

	msg := c.AddDomainAlias("test", "test.test")

	if _, ok := msg.(*common.ErrMsg); ok {
		t.Error("Expected to add domain alias, got: ", msg.GetText())
		return
	}

	msg = c.AddDomainAlias("test", "test.tst")

	if _, ok := msg.(*common.ErrMsg); ok {
		t.Error("Expected to add domain alias, got: ", msg.GetText())
		return
	}

	// "Refetch" the projects from the config file
	c.FetchProjects()

	c.RemoveDomainAlias("test", "test.test")

	// "Refetch" the projects from the config file
	c.FetchProjects()

	exists, project := c.ProjectExists("test")

	if !exists {
		t.Error("Expected project to exist")
		return
	}

	if len(project.DomainAliases) != 1 {
		t.Error("Expected 1 domain alias, got", len(project.DomainAliases))
		return
	}

	if project.DomainAliases[0].Value != "test.tst" {
		t.Error("Expected domain alias to be 'test.tst', got", project.DomainAliases[0])
	}

	// Check if the domain alias is removed from the hosts file
	hostsFilePath := c.GetConfig().GetHostsFile()
	hostsFile, err := os.Open(hostsFilePath)

	if err != nil {
		t.Error("Error to open hosts file, got:", err)
		return
	}

	defer hostsFile.Close()

	fileInfo, err := os.Stat(hostsFilePath)

	if err != nil {
		t.Error("Expected to stat hosts file, got", err)
		return
	}

	buf := make([]byte, fileInfo.Size())

	_, err = hostsFile.Read(buf)

	if err != nil {
		t.Error("Expected to read hosts file, got", err)
		return
	}

	hostsContent := string(buf)
	expected := "\n\n" +
		config.HostsBeginMarker +
		"\n127.0.0.1\ttest.local" +
		"\n127.0.0.1\ttest.tst" +
		config.HostsEndMarker

	if common.IsWindows() {
		hostsContent = strings.ReplaceAll(hostsContent, "\u0000", "")
	}

	if hostsContent != expected {
		t.Error("Expected hosts file to contain", expected, "got", hostsContent)
		return
	}
}
