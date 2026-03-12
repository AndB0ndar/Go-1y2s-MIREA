package storage

import (
	"sync"
	"time"
)

type ProcessedStore struct {
	mu    sync.RWMutex
	items map[string]time.Time
	ttl   time.Duration
}

func NewProcessedStore() *ProcessedStore {
	return &ProcessedStore{
		items: make(map[string]time.Time),
		ttl:   5 * time.Minute,
	}
}

func (s *ProcessedStore) Add(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[id] = time.Now().Add(s.ttl)
}

func (s *ProcessedStore) Exists(id string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	exp, ok := s.items[id]
	if !ok {
		return false
	}
	if time.Now().After(exp) {
		return false
	}
	return true
}
