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
	aIDs := make([]*account.ID, 0)
	for k := range s.T.Accounts {
		kCopy := k
		aIDs = append(aIDs, &kCopy)
	}
	return aIDs
}

func (s *S) CurrentProjects(aID *account.ID) []*project.ID {
	pIDs := make([]*project.ID, 0)
	if a, ok := s.T.Accounts[*aID]; ok {
		for k := range a.Projects {
			kCopy := k
			pIDs = append(pIDs, &kCopy)
		}
	}
	return pIDs
}

func (s *S) CurrentResources(aID *account.ID, pID *project.ID) []*resource.ID {
	rIDs := make([]*resource.ID, 0)
	if a, ok := s.T.Accounts[*aID]; ok {
		if p, ok := a.Projects[*pID]; ok {
			for k := range p.Resources {
				kCopy := k
				rIDs = append(rIDs, &kCopy)
			}
		}
	}
	return rIDs
}

func (s *S) CurrentAccountNames() []*account.Name {
	aNames := make([]*account.Name, 0)
	for aID := range s.T.Accounts {
		aNames = append(aNames, s.T.Accounts[aID].Name)
	}
	return aNames
}

func (s *S) CurrentProjectNames(aID *account.ID) []*project.Name {
	pNames := make([]*project.Name, 0)
	if a, ok := s.T.Accounts[*aID]; ok {
		for pID := range a.Projects {
			pNames = append(pNames, s.T.Accounts[*aID].Projects[pID].Name)
		}
	}
	return pNames
}

func (s *S) CurrentResourceNames(aID *account.ID, pID *project.ID) []*resource.Name {
	rNames := make([]*resource.Name, 0)
	if aTopo, ok := s.T.Accounts[*aID]; ok {
		if pTopo, ok := aTopo.Projects[*pID]; ok {
			for rID := range pTopo.Resources {
				rNames = append(rNames, pTopo.Resources[rID].Name)
			}
		}
	}
	return rNames
}

// GET

func (s *S) GetAccountID(aName *account.Name) *account.ID {
	for aID := range s.T.Accounts {
		if *s.T.Accounts[aID].Name == *aName {
			return &aID
		}
	}
	return &account.NotFoundID
}

func (s *S) GetProjectID(aID *account.ID, pName *project.Name) *project.ID {
	for pID := range s.T.Accounts[*aID].Projects {
		if *s.T.Accounts[*aID].Projects[pID].Name == *pName {
			return &pID
		}
	}
	return &project.NotFoundID
}

func (s *S) GetResourceID(aID *account.ID, pID *project.ID, rName *resource.Name) *resource.ID {
	for rID := range s.T.Accounts[*aID].Projects[*pID].Resources {
		if *s.T.Accounts[*aID].Projects[*pID].Resources[rID].Name == *rName {
			return &rID
		}
	}
	return &resource.NotFoundID
}
