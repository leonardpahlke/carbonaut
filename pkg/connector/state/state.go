package state

import (
	"log/slog"
	"sync"

	"carbonaut.dev/pkg/provider/resource"
	"carbonaut.dev/pkg/provider/topology"
)

type S struct {
	mutex sync.Mutex
	T     topology.T
}

func New() *S {
	var aIDCounter int32 = 0
	return &S{
		mutex: sync.Mutex{},
		T: topology.T{
			Accounts:          make(map[topology.AccountID]*topology.AccountT),
			AccountsIDCounter: &aIDCounter,
		},
	}
}

// ADD / CREATE

func (s *S) AddAccount(a *topology.AccountT) *topology.AccountID {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	*s.T.AccountsIDCounter += 1
	aID := topology.AccountID(*s.T.AccountsIDCounter)
	s.T.Accounts[aID] = a
	return &aID
}

func (s *S) AddProject(aID *topology.AccountID, p *topology.ProjectT) *topology.ProjectID {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	*s.T.Accounts[*aID].ProjectIDCounter += 1
	pID := topology.ProjectID(*s.T.Accounts[*aID].ProjectIDCounter)
	s.T.Accounts[*aID].Projects[pID] = p
	return &pID
}

func (s *S) AddResource(aID *topology.AccountID, pID *topology.ProjectID, r *topology.ResourceT) *topology.ResourceID {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	*s.T.Accounts[*aID].Projects[*pID].ResourceIDCounter += 1
	rID := topology.ResourceID(*s.T.Accounts[*aID].Projects[*pID].ResourceIDCounter)
	s.T.Accounts[*aID].Projects[*pID].Resources[rID] = r
	return &rID
}

// DELETE

func (s *S) RemoveAccount(aID *topology.AccountID) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.T.Accounts, *aID)
}

func (s *S) RemoveProject(aID *topology.AccountID, pID *topology.ProjectID) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.T.Accounts[*aID].Projects, *pID)
}

func (s *S) RemoveProjects(aID *topology.AccountID, pIDs []*topology.ProjectID) {
	if len(pIDs) == 0 {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i := range pIDs {
		slog.Info("delete project", "project ID", pIDs[i])
		delete(s.T.Accounts[*aID].Projects, *pIDs[i])
	}
}

func (s *S) RemoveProjectsByName(aID *topology.AccountID, pNames []*resource.ProjectName) {
	if len(pNames) == 0 {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i := range pNames {
		slog.Info("delete project", "project name", *pNames[i])
		delete(s.T.Accounts[*aID].Projects, *s.GetProjectID(aID, pNames[i]))
	}
}

func (s *S) RemoveResources(aID *topology.AccountID, pID *topology.ProjectID, rIDs []*topology.ResourceID) {
	if len(rIDs) == 0 {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i := range rIDs {
		slog.Info("delete resource", "resource name", *rIDs[i])
		delete(s.T.Accounts[*aID].Projects[*pID].Resources, *rIDs[i])
	}
}

func (s *S) RemoveResourceByName(aID *topology.AccountID, pID *topology.ProjectID, rNames []*resource.ResourceName) {
	if len(rNames) == 0 {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i := range rNames {
		slog.Info("delete resource", "resource name", *rNames[i])
		delete(s.T.Accounts[*aID].Projects[*pID].Resources, *s.GetResourceID(aID, pID, rNames[i]))
	}
}

func (s *S) RemoveResource(aID *topology.AccountID, pID *topology.ProjectID, rID *topology.ResourceID) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	slog.Info("delete resource", "resource name", *rID)
	delete(s.T.Accounts[*aID].Projects[*pID].Resources, *rID)
}

// COLLECT

func (s *S) CurrentAccounts() []*topology.AccountID {
	aIDs := make([]*topology.AccountID, 0)
	for k := range s.T.Accounts {
		kCopy := k
		aIDs = append(aIDs, &kCopy)
	}
	return aIDs
}

func (s *S) CurrentProjects(aID *topology.AccountID) []*topology.ProjectID {
	pIDs := make([]*topology.ProjectID, 0)
	if a, ok := s.T.Accounts[*aID]; ok {
		for k := range a.Projects {
			kCopy := k
			pIDs = append(pIDs, &kCopy)
		}
	}
	return pIDs
}

func (s *S) CurrentResources(aID *topology.AccountID, pID *topology.ProjectID) []*topology.ResourceID {
	rIDs := make([]*topology.ResourceID, 0)
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

func (s *S) CurrentAccountNames() []*resource.AccountName {
	aNames := make([]*resource.AccountName, 0)
	for aID := range s.T.Accounts {
		aNames = append(aNames, s.T.Accounts[aID].Name)
	}
	return aNames
}

func (s *S) CurrentProjectNames(aID *topology.AccountID) []*resource.ProjectName {
	pNames := make([]*resource.ProjectName, 0)
	if a, ok := s.T.Accounts[*aID]; ok {
		for pID := range a.Projects {
			pNames = append(pNames, s.T.Accounts[*aID].Projects[pID].Name)
		}
	}
	return pNames
}

func (s *S) CurrentResourceNames(aID *topology.AccountID, pID *topology.ProjectID) []*resource.ResourceName {
	rNames := make([]*resource.ResourceName, 0)
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

func (s *S) GetAccountID(aName *resource.AccountName) *topology.AccountID {
	for aID := range s.T.Accounts {
		if *s.T.Accounts[aID].Name == *aName {
			return &aID
		}
	}
	return &topology.AccountNotFoundID
}

func (s *S) GetProjectID(aID *topology.AccountID, pName *resource.ProjectName) *topology.ProjectID {
	for pID := range s.T.Accounts[*aID].Projects {
		if *s.T.Accounts[*aID].Projects[pID].Name == *pName {
			return &pID
		}
	}
	return &topology.ProjectNotFoundID
}

func (s *S) GetResourceID(aID *topology.AccountID, pID *topology.ProjectID, rName *resource.ResourceName) *topology.ResourceID {
	for rID := range s.T.Accounts[*aID].Projects[*pID].Resources {
		if *s.T.Accounts[*aID].Projects[*pID].Resources[rID].Name == *rName {
			return &rID
		}
	}
	return &topology.ResourceNotFoundID
}
