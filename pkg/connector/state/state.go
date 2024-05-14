package state

import (
	"sync"

	"carbonaut.dev/pkg/schema/provider"
	"carbonaut.dev/pkg/schema/provider/data/account"
	"carbonaut.dev/pkg/schema/provider/data/account/project"
	"carbonaut.dev/pkg/schema/provider/data/account/project/resource"
)

type S struct {
	mutex sync.Mutex
	T     provider.Topology
}

func New() *S {
	var aIDCounter int32 = 0
	return &S{
		mutex: sync.Mutex{},
		T: provider.Topology{
			Accounts:          make(map[account.ID]*account.Topology),
			AccountsIDCounter: &aIDCounter,
		},
	}
}

// ADD / CREATE

func (s *S) AddAccount(a *account.Topology) *account.ID {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	*s.T.AccountsIDCounter += 1
	aID := account.ID(*s.T.AccountsIDCounter)
	s.T.Accounts[aID] = a
	return &aID
}

func (s *S) AddProject(aID *account.ID, p *project.Topology) *project.ID {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	*s.T.Accounts[*aID].ProjectIDCounter += 1
	pID := project.ID(*s.T.Accounts[*aID].ProjectIDCounter)
	s.T.Accounts[*aID].Projects[pID] = p
	return &pID
}

func (s *S) AddResource(aID *account.ID, pID *project.ID, r *resource.Topology) *resource.ID {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	*s.T.Accounts[*aID].Projects[*pID].ResourceIDCounter += 1
	rID := resource.ID(*s.T.Accounts[*aID].Projects[*pID].ResourceIDCounter)
	s.T.Accounts[*aID].Projects[*pID].Resources[rID] = r
	return &rID
}

// func (s *S) AddProjectResources(aID *account.ID, pID *project.ID, r project.Resources) {
// 	s.mutex.Lock()
// 	defer s.mutex.Unlock()
// 	s.T.Accounts[*aID].Projects[*pID] = &project.Topology{
// 		Resources: r,
// 		CreatedAt: time.Now(),
// 	}
// }

// DELETE

func (s *S) RemoveAccount(aID *account.ID) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.T.Accounts, *aID)
}

func (s *S) RemoveProject(aID *account.ID, pID *project.ID) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.T.Accounts[*aID].Projects, *pID)
}

func (s *S) RemoveProjects(aID *account.ID, pIDs []*project.ID) {
	if len(pIDs) == 0 {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i := range pIDs {
		delete(s.T.Accounts[*aID].Projects, *pIDs[i])
	}
}

func (s *S) RemoveProjectsByName(aID *account.ID, pIDs []*project.Name) {
	if len(pIDs) == 0 {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i := range pIDs {
		delete(s.T.Accounts[*aID].Projects, *s.GetProjectID(aID, pIDs[i]))
	}
}

func (s *S) RemoveResources(aID *account.ID, pID *project.ID, rIDs []*resource.ID) {
	if len(rIDs) == 0 {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i := range rIDs {
		delete(s.T.Accounts[*aID].Projects[*pID].Resources, *rIDs[i])
	}
}

func (s *S) RemoveResourceByName(aID *account.ID, pID *project.ID, rIDs []*resource.Name) {
	if len(rIDs) == 0 {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i := range rIDs {
		delete(s.T.Accounts[*aID].Projects[*pID].Resources, *s.GetResourceID(aID, pID, rIDs[i]))
	}
}

func (s *S) RemoveResource(aID *account.ID, pID *project.ID, rID *resource.ID) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.T.Accounts[*aID].Projects[*pID].Resources, *rID)
}

// COLLECT

func (s *S) CurrentAccounts() []*account.ID {
	currentAccountIDs := make([]*account.ID, 0)
	for k := range s.T.Accounts {
		currentAccountIDs = append(currentAccountIDs, &k)
	}
	return currentAccountIDs
}

func (s *S) CurrentProjects(aID *account.ID) []*project.ID {
	currentProjectIDs := make([]*project.ID, 0)
	if account, ok := s.T.Accounts[*aID]; ok {
		for k := range account.Projects {
			currentProjectIDs = append(currentProjectIDs, &k)
		}
	}
	return currentProjectIDs
}

func (s *S) CurrentResources(aID *account.ID, pID *project.ID) []*resource.ID {
	currentResourceIDs := make([]*resource.ID, 0)
	if account, ok := s.T.Accounts[*aID]; ok {
		if project, ok := account.Projects[*pID]; ok {
			for k := range project.Resources {
				currentResourceIDs = append(currentResourceIDs, &k)
			}
		}
	}
	return currentResourceIDs
}

func (s *S) CurrentAccountNames() []*account.Name {
	accountNames := make([]*account.Name, 0)
	for aID := range s.T.Accounts {
		accountNames = append(accountNames, s.T.Accounts[aID].Name)
	}
	return accountNames
}

func (s *S) CurrentProjectNames(aID *account.ID) []*project.Name {
	projectNames := make([]*project.Name, 0)
	if account, ok := s.T.Accounts[*aID]; ok {
		for pID := range account.Projects {
			projectNames = append(projectNames, s.T.Accounts[*aID].Projects[pID].Name)
		}
	}
	return projectNames
}

func (s *S) CurrentResourceNames(aID *account.ID, pID *project.ID) []*resource.Name {
	resourceNames := make([]*resource.Name, 0)
	if aTopo, ok := s.T.Accounts[*aID]; ok {
		if pTopo, ok := aTopo.Projects[*pID]; ok {
			for rID := range pTopo.Resources {
				resourceNames = append(resourceNames, pTopo.Resources[rID].Name)
			}
		}
	}
	return resourceNames
}

// GET

func (s *S) GetAccountID(aName *account.Name) *account.ID {
	for i := range s.T.Accounts {
		if *s.T.Accounts[i].Name == *aName {
			return &i
		}
	}
	return &account.NotFoundID
}

func (s *S) GetProjectID(aID *account.ID, pName *project.Name) *project.ID {
	for i := range s.T.Accounts[*aID].Projects {
		if *s.T.Accounts[*aID].Projects[i].Name == *pName {
			return &i
		}
	}
	return &project.NotFoundID
}

func (s *S) GetResourceID(aID *account.ID, pID *project.ID, rName *resource.Name) *resource.ID {
	for i := range s.T.Accounts[*aID].Projects[*pID].Resources {
		if *s.T.Accounts[*aID].Projects[*pID].Resources[i].Name == *rName {
			return &i
		}
	}
	return &resource.NotFoundID
}
