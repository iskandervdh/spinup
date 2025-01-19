package core

import (
	"testing"

	"github.com/iskandervdh/spinup/common"
)

func TestAddDomainAlias(t *testing.T) {
	c := TestingCore("add_domain_alias")

	// Fetch the commands and projects from their config files
	c.FetchCommands()
	c.FetchProjects()

	c.AddProject("test", 1234, []string{})

	// "Refetch" the projects from the config file
	c.FetchProjects()

	c.AddDomainAlias("test", "test.test.test")

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

	if project.DomainAliases[0].Value != "test.test.test" {
		t.Error("Expected domain alias to be 'test.test.test', got", project.DomainAliases[0])
	}
}

func TestRemoveDomainAlias(t *testing.T) {
	c := TestingCore("remove_domain_alias")

	// Fetch the commands and projects from their config files
	c.FetchCommands()
	c.FetchProjects()

	c.AddProject("test", 1234, []string{})

	err := c.FetchProjects()

	if err != nil {
		t.Error("Expected to fetch projects, got", err)
		return
	}

	msg := c.AddDomainAlias("test", "test.test.test")

	if _, ok := msg.(*common.ErrMsg); ok {
		t.Error("Expected to add domain alias, got: ", msg.GetText())
		return
	}

	msg = c.AddDomainAlias("test", "tst.test.test")

	if _, ok := msg.(*common.ErrMsg); ok {
		t.Error("Expected to add domain alias, got: ", msg.GetText())
		return
	}

	// "Refetch" the projects from the config file
	c.FetchProjects()

	c.RemoveDomainAlias("test", "test.test.test")

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

	if project.DomainAliases[0].Value != "tst.test.test" {
		t.Error("Expected domain alias to be 'tst.test.test', got", project.DomainAliases[0])
	}
}
