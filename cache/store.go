// Package cache provides a thread-safe in-memory LRU cache for storing
// receipt data and calculated points. It uses a read-write mutex to
// handle concurrent access safely. Items are evicted when the cache
// reaches its size limit, ensuring efficient access to frequently used data.
package cache

import (
	"log"
	"sync"

	lru "github.com/hashicorp/golang-lru"
)

// Store wraps an LRU cache with a mutex for safe concurrent access.
type Store[T any] struct {
	mu    sync.RWMutex
	cache *lru.Cache
}

// NewStore creates a new LRU cache with the given size.
func NewStore[T any](size int) *Store[T] {
	c, err := lru.New(size)
	if err != nil {
		log.Fatalf("failed to create LRU cache: %v", err)
	}
	return &Store[T]{cache: c}
}

// Get retrieves a value from the cache by key.
func (s *Store[T]) Get(key string) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.cache.Get(key)
	if !ok {
		var zero T
		return zero, false
	}
	v, ok := value.(T)
	if !ok {
		var zero T
		return zero, false
	}
	return v, true
}

// Set stores a key-value pair in the cache.
func (s *Store[T]) Set(key string, value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache.Add(key, value)
}
