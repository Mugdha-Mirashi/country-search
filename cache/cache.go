package cache

import (
	"sync"
)

// Cache interface defines the methods for the cache
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
}

// inMemoryCache is a simple thread-safe in-memory cache
type inMemoryCache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

// NewCache creates a new instance of inMemoryCache
func NewCache() Cache {
	return &inMemoryCache{
		data: make(map[string]interface{}),
	}
}

// Get retrieves a value from the cache
func (c *inMemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, found := c.data[key]
	return value, found
}

// Set stores a value in the cache
func (c *inMemoryCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}