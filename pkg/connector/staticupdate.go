package connector

import (
	"fmt"
	"time"

	"carbonaut.dev/pkg/plugins/staticresplugins"
	"carbonaut.dev/pkg/schema/provider/data/account"
	"carbonaut.dev/pkg/schema/provider/data/account/project"
	"carbonaut.dev/pkg/schema/provider/data/account/project/resource"
	"carbonaut.dev/pkg/schema/provider/types/staticres"
	"carbonaut.dev/pkg/util/compareutils"
)

func (c *C) updateStaticData(aID *account.ID) error {
	var pRes staticres.Provider
	if rConfig, ok := (c.providerConfig.Resources)[*aID]; ok {
		if p, found := staticresplugins.GetPlugin(rConfig.StaticResConfig.Plugin); found {
			pRes = p
		} else {
			return fmt.Errorf("could not find plugin %s for account %s", *rConfig.StaticResConfig.Plugin, *aID)
		}
	}

	discoveredProjectIDs, err := c.fetchProjects(aID, pRes)
	if err != nil {
		return err
	}

	remainingProjects, toBeDeletedProjects, toBeCreatedProjects := compareutils.CompareLists(discoveredProjectIDs, c.state.CurrentProjects(aID))

	c.log.Debug("fetched projects", "account", *aID, "remainingProjects", remainingProjects, "toBeDeletedProjects", toBeDeletedProjects, "toBeCreatedProjects", toBeCreatedProjects)

	// INFO: remove all resources that are not found anymore but loaded to state
	c.state.RemoveProjects(aID, toBeDeletedProjects)

	// INFO: fetch resources from remainingProjects since resources may change
	if err := c.updateProjectResources(aID, remainingProjects, pRes); err != nil {
		return err
	}

	if err := c.updateProjectResources(aID, toBeCreatedProjects, pRes); err != nil {
		return err
	}

	return nil
}

func (c *C) fetchProjects(aID *account.ID, p staticres.Provider) ([]*project.ID, error) {
	var projects []*project.ID
	if rConfig, ok := (c.providerConfig.Resources)[*aID]; ok {
		discoveredProjects, err := p.DiscoverProjectIdentifiers(rConfig.StaticResConfig)
		if err != nil {
			return nil, fmt.Errorf("unable to list resources for account: %s, error: %v", *aID, err)
		}
		projects = discoveredProjects
	} else {
		return nil, fmt.Errorf("resource config not found in account: %s", *aID)
	}
	return projects, nil
}

func (c *C) fetchResources(aID *account.ID, pID *project.ID, p staticres.Provider) ([]*resource.ID, error) {
	var r []*resource.ID
	if rConfig, ok := (c.providerConfig.Resources)[*aID]; ok {
		discoveredResources, err := p.DiscoverStaticResourceIdentifiers(rConfig.StaticResConfig, pID)
		if err != nil {
			return nil, fmt.Errorf("unable to list resources for account: %s for project: %s, error: %v", *aID, *pID, err)
		}
		r = discoveredResources
	} else {
		return nil, fmt.Errorf("resource config not found in account: %s", *aID)
	}
	return r, nil
}

func (c *C) getStaticResourcesData(aID *account.ID, pID *project.ID, rIDs []*resource.ID, p staticres.Provider) (project.Resources, error) {
	d := make(project.Resources)
	for i := range rIDs {
		if rConfig, ok := (c.providerConfig.Resources)[*aID]; ok {
			staticResData, err := p.GetStaticResourceData(rConfig.StaticResConfig, pID, rIDs[i])
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve resource information for account: %s in project: %s for resource: %s, error: %v", *aID, *pID, *rIDs[i], err)
			}
			d[*rIDs[i]] = &resource.Topology{
				StaticData: staticResData,
				CreatedAt:  time.Now(),
				Plugin:     p.GetName(),
			}
		} else {
			c.log.Error("resource config not found in account", "account", aID)
			return nil, fmt.Errorf("resource config not found in account: %s, project: %s, resource:%s", *aID, *pID, *rIDs[i])
		}
	}
	return d, nil
}

func (c *C) updateProjectResources(aID *account.ID, discoveredProjects []*project.ID, pRes staticres.Provider) error {
	for i := range discoveredProjects {
		discoveredResourceIDs, err := c.fetchResources(aID, discoveredProjects[i], pRes)
		if err != nil {
			return fmt.Errorf("could not fetchResources for project %s in account %s, error: %v", *discoveredProjects[i], *aID, err)
		}
		remainingResources, toBeDeletedResources, toBeCreatedResources := compareutils.CompareLists(discoveredResourceIDs, c.state.CurrentResources(aID, discoveredProjects[i]))

		c.log.Debug("fetched resources", "account", *aID, "project", *discoveredProjects[i], "remainingResources", remainingResources, "toBeDeletedResources", toBeDeletedResources, "toBeCreatedResources", toBeCreatedResources)

		// INFO: do nothing about the remainingResources since they are already loaded to state

		// INFO: remove all resources that are not found anymore but loaded to state
		c.state.RemoveResources(aID, discoveredProjects[i], toBeDeletedResources)

		// INFO: add all new resources that were found but are not loaded to state yet
		resourceData, err := c.getStaticResourcesData(aID, discoveredProjects[i], toBeCreatedResources, pRes)
		if err != nil {
			return fmt.Errorf("could not get resource data for project %s in account %s, error: %v", *discoveredProjects[i], *aID, err)
		}
		c.state.AddProjectResources(aID, discoveredProjects[i], resourceData)
	}
	return nil
}
