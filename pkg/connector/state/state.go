package state

import (
	"sync"
	"time"

	"carbonaut.dev/pkg/schema/provider"
	"carbonaut.dev/pkg/schema/provider/data/account"
	"carbonaut.dev/pkg/schema/provider/data/account/project"
	"carbonaut.dev/pkg/schema/provider/data/account/project/resource"
)

type S struct {
	mutex    sync.Mutex
	Accounts provider.Topology
}

func New() *S {
	return &S{
		mutex:    sync.Mutex{},
		Accounts: make(provider.Topology),
	}
}

// ADD / CREATE

func (s *S) AddAccount(aID *account.ID, a *account.Topology) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Accounts[*aID] = a
}

func (s *S) AddProject(aID *account.ID, pID *project.ID, p *project.Topology) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Accounts[*aID].Projects[*pID] = p
}

func (s *S) AddResource(aID *account.ID, pID *project.ID, rID *resource.ID, r *resource.Topology) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Accounts[*aID].Projects[*pID].Resources[*rID] = r
}

func (s *S) AddProjectResources(aID *account.ID, pID *project.ID, r project.Resources) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Accounts[*aID].Projects[*pID] = &project.Topology{
		Resources: r,
		CreatedAt: time.Now(),
	}
}

// DELETE

func (s *S) RemoveAccount(aID *account.ID) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.Accounts, *aID)
}

func (s *S) RemoveProject(aID *account.ID, pID *project.ID) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.Accounts[*aID].Projects, *pID)
}

func (s *S) RemoveProjects(aID *account.ID, pIDs []*project.ID) {
	if len(pIDs) == 0 {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i := range pIDs {
		delete(s.Accounts[*aID].Projects, *pIDs[i])
	}
}

func (s *S) RemoveResources(aID *account.ID, pID *project.ID, rIDs []*resource.ID) {
	if len(rIDs) == 0 {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for i := range rIDs {
		delete(s.Accounts[*aID].Projects[*pID].Resources, *rIDs[i])
	}
}

func (s *S) RemoveResource(aID *account.ID, pID *project.ID, rID *resource.ID) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.Accounts[*aID].Projects[*pID].Resources, *rID)
}

// COLLECT

func (s *S) CurrentAccounts() []*account.ID {
	currentAccountIDs := make([]*account.ID, 0)
	for k := range s.Accounts {
		currentAccountIDs = append(currentAccountIDs, &k)
	}
	return currentAccountIDs
}

func (s *S) CurrentProjects(aID *account.ID) []*project.ID {
	currentProjectIDs := make([]*project.ID, 0)
	if account, ok := s.Accounts[*aID]; ok {
		for k := range account.Projects {
			currentProjectIDs = append(currentProjectIDs, &k)
		}
	}
	return currentProjectIDs
}

func (s *S) CurrentResources(aID *account.ID, pID *project.ID) []*resource.ID {
	currentResourceIDs := make([]*resource.ID, 0)
	if account, ok := s.Accounts[*aID]; ok {
		if project, ok := account.Projects[*pID]; ok {
			for k := range project.Resources {
				currentResourceIDs = append(currentResourceIDs, &k)
			}
		}
	}
	return currentResourceIDs
}
