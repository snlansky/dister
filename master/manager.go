package master

import (
	"sync"
)

type Manager struct {
	mu      sync.RWMutex
	workers map[string]*Worker
}

func NewManager() *Manager {
	return &Manager{
		workers: map[string]*Worker{},
	}
}

func (man *Manager) AddWorker(worker *Worker) {
	man.mu.Lock()
	defer man.mu.Unlock()
	man.workers[worker.id] = worker
}

func (man *Manager) DeleteWorker(id string) *Worker {
	man.mu.Lock()
	defer man.mu.Unlock()

	if w, ok := man.workers[id]; ok {
		delete(man.workers, id)
		return w
	} else {
		return nil
	}
}
