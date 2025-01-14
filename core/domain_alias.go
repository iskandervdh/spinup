package core

import (
	"fmt"

	"github.com/iskandervdh/spinup/common"
	"github.com/iskandervdh/spinup/database/sqlc"
)

type DomainAlias = sqlc.DomainAlias

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
	if common.GetDomain(project.Name) == domainAlias {
		return common.NewErrMsg("Domain alias '%s' is already the domain of project '%s'", domainAlias, projectName)
	}

	for projectName, project := range c.projects {
		// Check if the domain alias is the domain of another project
		if common.GetDomain(project.Name) == domainAlias {
			return common.NewErrMsg("Domain alias '%s' is already the domain of project '%s'", domainAlias, projectName)
		}

		// Check if the domain alias is already a domain alias of any project
		for _, alias := range project.DomainAliases {
			if alias.Value == domainAlias {
				return common.NewErrMsg("Domain alias '%s' already exists on project '%s'", domainAlias, projectName)
			}
		}
	}

	err := c.config.NginxAddDomainAlias(projectName, domainAlias)

	if err != nil {
		return common.NewErrMsg(fmt.Sprintln("Error trying to add domain alias to nginx config file", err))
	}

	err = c.dbQueries.CreateDomainAlias(c.dbContext, sqlc.CreateDomainAliasParams{
		Value:     domainAlias,
		ProjectID: project.ID,
	})

	if err != nil {
		return common.NewErrMsg("Error adding domain alias to database: %s", err)
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

	for i, alias := range project.DomainAliases {
		if alias.Value == domainAlias {
			// Remove the given domain alias from the domain aliases array of the project
			project.DomainAliases = append(project.DomainAliases[:i], project.DomainAliases[i+1:]...)

			err := c.dbQueries.DeleteDomainAlias(c.dbContext, sqlc.DeleteDomainAliasParams{
				Value:     domainAlias,
				ProjectID: project.ID,
			})

			if err != nil {
				return common.NewErrMsg("Error removing domain alias from database: %s", err)
			}

			return common.NewSuccessMsg("Removed domain alias '%s' from project '%s'", domainAlias, projectName)
		}
	}

	return common.NewErrMsg("Domain alias '%s' does not exist on project '%s'", domainAlias, projectName)
}
