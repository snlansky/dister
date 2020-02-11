package master

import (
	"math/rand"
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

func (man *Manager) GetRandWorker() *Worker {
	man.mu.Lock()
	defer man.mu.Unlock()

	index := rand.Intn(len(man.workers))
	i := 0
	for _, w := range man.workers {
		if index == i {
			return w
		}
		i++
	}
	return nil
}
