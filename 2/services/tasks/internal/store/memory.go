package store

import (
	"github.com/google/uuid"
	"sync"
)

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	Done        bool   `json:"done"`
}

type MemoryStore struct {
	mu    sync.RWMutex
	tasks map[string]Task
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks: make(map[string]Task),
	}
}

func (s *MemoryStore) Create(task Task) Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	task.ID = uuid.New().String()[:8]
	s.tasks[task.ID] = task
	return task
}

func (s *MemoryStore) List() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) Get(id string) (Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	return t, ok
}

func (s *MemoryStore) Update(id string, updated Task) (Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.tasks[id]
	if !ok {
		return Task{}, false
	}
	if updated.Title != "" {
		existing.Title = updated.Title
	}
	if updated.Description != "" {
		existing.Description = updated.Description
	}
	if updated.DueDate != "" {
		existing.DueDate = updated.DueDate
	}
	existing.Done = updated.Done
	s.tasks[id] = existing
	return existing, true
}

func (s *MemoryStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; ok {
		delete(s.tasks, id)
		return true
	}
	return false
}
