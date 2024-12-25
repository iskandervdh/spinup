package core

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/iskandervdh/spinup/common"
)

// DomainAliases is a map of domain aliases to their domain names.
type DomainAliases []string

// Add a domain alias to the given project.
func (c *Core) AddDomainAlias(projectName string, domainAlias string) common.Msg {
	c.RequireSudo()

	if c.projects == nil {
		return common.NewErrMsg("No projects found")
	}

	exists, project := c.ProjectExists(projectName)

	if !exists {
		return common.NewErrMsg("Project '%s' does not exist", projectName)
	}

	// Check if the domain alias is already defined as the domain of the project
	if project.Domain == domainAlias {
		return common.NewErrMsg("Domain alias '%s' is already the domain of project '%s'", domainAlias, projectName)
	}

	for projectName, project := range c.projects {
		// Check if the domain alias is the domain of another project
		if project.Domain == domainAlias {
			return common.NewErrMsg("Domain alias '%s' is already the domain of project '%s'", domainAlias, projectName)
		}

		// Check if the domain alias is already a domain alias of any project
		for _, alias := range project.DomainAliases {
			if alias == domainAlias {
				return common.NewErrMsg("Domain alias '%s' already exists on project '%s'", domainAlias, projectName)
			}
		}
	}

	err := c.config.NginxAddDomainAlias(projectName, domainAlias)

	if err != nil {
		return common.NewErrMsg(fmt.Sprintln("Error trying to add domain alias to nginx config file", err))
	}

	err = c.config.AddDomain(domainAlias)

	if err != nil {
		// Remove the domain alias from the nginx config file if adding domain alias to hosts file fails
		c.config.NginxRemoveDomainAlias(projectName, domainAlias)

		return common.NewErrMsg(fmt.Sprintln("Error trying to add domain to hosts file", err))
	}

	project.DomainAliases = append(project.DomainAliases, domainAlias)

	c.projects[projectName] = project

	updatedProjectsConfig, err := json.MarshalIndent(c.projects, "", "  ")

	if err != nil {
		return common.NewErrMsg("Error encoding projects to json: %s", err)
	}

	err = os.WriteFile(c.getProjectsFilePath(), updatedProjectsConfig, 0644)

	if err != nil {
		return common.NewErrMsg("Error writing projects to file: %s", err)
	}

	return common.NewSuccessMsg("Added domain alias '%s' to project '%s'", domainAlias, projectName)
}

// Remove a domain alias from the given project.
func (c *Core) RemoveDomainAlias(projectName string, domainAlias string) common.Msg {
	if c.projects == nil {
		return common.NewErrMsg("No projects found")
	}

	exists, project := c.ProjectExists(projectName)

	if !exists {
		return common.NewErrMsg("Project '%s' does not exist", projectName)
	}

	err := c.config.NginxRemoveDomainAlias(projectName, domainAlias)

	if err != nil {
		return common.NewErrMsg(fmt.Sprintln("Error trying to remove domain alias from nginx config file", err))
	}

	err = c.config.RemoveDomain(domainAlias)

	if err != nil {
		// Readd the domain alias from the nginx config file if removing the domain alias from hosts file fails
		c.config.NginxRemoveDomainAlias(projectName, domainAlias)

		return common.NewErrMsg(fmt.Sprintln("Error trying to add domain to hosts file", err))
	}

	for i, alias := range project.DomainAliases {
		if alias == domainAlias {
			// Remove the given domain alias from the domain aliases array of the project
			project.DomainAliases = append(project.DomainAliases[:i], project.DomainAliases[i+1:]...)

			c.projects[projectName] = project

			updatedProjectsConfig, err := json.MarshalIndent(c.projects, "", "  ")

			if err != nil {
				return common.NewErrMsg("Error encoding projects to json: %s", err)
			}

			err = os.WriteFile(c.getProjectsFilePath(), updatedProjectsConfig, 0644)

			if err != nil {
				return common.NewErrMsg("Error writing projects to file: %s", err)
			}

			return common.NewSuccessMsg("Removed domain alias '%s' from project '%s'", domainAlias, projectName)
		}
	}

	return common.NewErrMsg("Domain alias '%s' does not exist on project '%s'", domainAlias, projectName)
}
