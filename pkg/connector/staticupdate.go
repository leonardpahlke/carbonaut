package connector

import (
	"fmt"
	"log/slog"
	"time"

	"carbonaut.dev/pkg/plugin/staticresplugins"
	"carbonaut.dev/pkg/provider/data/account"
	"carbonaut.dev/pkg/provider/data/account/project"
	"carbonaut.dev/pkg/provider/data/account/project/resource"
	"carbonaut.dev/pkg/provider/types/staticres"
	"carbonaut.dev/pkg/util/compareutils"
	"go.uber.org/multierr"
)

func (c *C) updateStaticData(aID *account.ID) error {
	var pRes staticres.Provider
	pRes, err := staticresplugins.GetPlugin(c.state.T.Accounts[*aID].Config)
	if err != nil {
		return fmt.Errorf("could not find plugin %v for account %v, error: %v", *c.state.T.Accounts[*aID].Config.Plugin, *aID, err)
	}

	discoveredProjectIDs, err := pRes.DiscoverProjectIdentifiers()
	if err != nil {
		return err
	}

	remainingProjects, toBeDeletedProjects, toBeCreatedProjects := compareutils.CompareLists(discoveredProjectIDs, c.state.CurrentProjectNames(aID))

	slog.Debug("fetched projects", "account", *aID, "remainingProjects", toStringSliceProjects(remainingProjects), "toBeDeletedProjects", toStringSliceProjects(toBeDeletedProjects), "toBeCreatedProjects", toStringSliceProjects(toBeCreatedProjects))

	// INFO: remove all resources that are not found anymore but loaded to state
	c.state.RemoveProjectsByName(aID, toBeDeletedProjects)

	remainingProjectIDs := make([]*project.ID, 0)
	toBeCreatedProjectIDs := make([]*project.ID, 0)

	for i := range remainingProjects {
		remainingProjectIDs = append(remainingProjectIDs, c.state.GetProjectID(aID, remainingProjects[i]))
	}

	for i := range toBeCreatedProjects {
		toBeCreatedProjectIDs = append(toBeCreatedProjectIDs, c.state.AddProject(aID, &project.Topology{
			Name:              toBeCreatedProjects[i],
			Resources:         make(map[resource.ID]*resource.Topology),
			CreatedAt:         time.Now(),
			ResourceIDCounter: new(int32),
		}))
	}

	// INFO: fetch resources from remainingProjects since resources may change
	if err := c.updateProjectResources(aID, remainingProjectIDs, pRes); err != nil {
		return err
	}

	if err := c.updateProjectResources(aID, toBeCreatedProjectIDs, pRes); err != nil {
		return err
	}

	return nil
}

func (c *C) updateProjectResources(aID *account.ID, pIDs []*project.ID, pRes staticres.Provider) error {
	var mError error
	for i := range pIDs {
		discoveredResourceNames, err := pRes.DiscoverStaticResourceIdentifiers(c.state.T.Accounts[*aID].Projects[*pIDs[i]].Name)
		if err != nil {
			mError = multierr.Append(mError, fmt.Errorf("could not fetchResources for project %v in account %v, error: %v", *pIDs[i], *aID, err))
			continue
		}
		remainingResourceNames, toBeDeletedResourceNames, toBeCreatedResourceNames := compareutils.CompareLists(discoveredResourceNames, c.state.CurrentResourceNames(aID, pIDs[i]))

		slog.Debug("fetched resources", "account", *aID, "project", *pIDs[i], "remainingResources", toStringSliceResources(remainingResourceNames), "toBeDeletedResources", toStringSliceResources(toBeDeletedResourceNames), "toBeCreatedResources", toStringSliceResources(toBeCreatedResourceNames))

		// INFO: do nothing about the remainingResources since they are already loaded to state

		// INFO: remove all resources that are not found anymore but loaded to state
		c.state.RemoveResourceByName(aID, pIDs[i], toBeDeletedResourceNames)

		for j := range toBeCreatedResourceNames {
			data, err := pRes.GetStaticResourceData(c.state.T.Accounts[*aID].Projects[*pIDs[i]].Name, toBeCreatedResourceNames[j])
			if err != nil {
				mError = multierr.Append(mError, fmt.Errorf("could not GetStaticResourceData for resource %v in project %v in account %v, error: %v", *toBeCreatedResourceNames[j], *pIDs[i], *aID, err))
				continue
			}
			c.state.AddResource(aID, pIDs[i], &resource.Topology{
				Name:       toBeCreatedResourceNames[j],
				StaticData: data,
				CreatedAt:  time.Now(),
				Plugin:     pRes.GetName(),
			})
		}
	}
	return mError
}

// toStringSliceProjects converts a slice of pointers to project.Name to a slice of strings
func toStringSliceProjects(names []*project.Name) []string {
	var result []string
	for _, name := range names {
		if name != nil {
			result = append(result, string(*name))
		} else {
			result = append(result, "nil")
		}
	}
	return result
}

// toStringSliceResources converts a slice of pointers to resource.Name to a slice of strings
func toStringSliceResources(names []*resource.Name) []string {
	var result []string
	for _, name := range names {
		if name != nil {
			result = append(result, string(*name))
		} else {
			result = append(result, "nil")
		}
	}
	return result
}
