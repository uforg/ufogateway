package cache

import (
	"sync"
	"time"
)

// cacheItem represents a cache cacheItem with a value and an expiration time.
type cacheItem struct {
	value  any
	expiry time.Time
}

// isExpired checks if the cache item has expired.
func (i cacheItem) isExpired() bool {
	return time.Now().After(i.expiry)
}

// CacheInstance is a generic cache implementation with support for time-to-live
// (TTL) expiration.
type CacheInstance struct {
	items map[string]cacheItem // The map storing cache items.
	mu    sync.Mutex           // Mutex for controlling concurrent access to the cache.
}

// NewCacheInstance creates a new Cache instance and starts a goroutine to periodically
// remove expired items every 10 seconds.
func NewCacheInstance() *CacheInstance {
	c := &CacheInstance{
		items: make(map[string]cacheItem),
	}

	go func() {
		for range time.Tick(10 * time.Second) {
			c.mu.Lock()

			for key, item := range c.items {
				if item.isExpired() {
					delete(c.items, key)
				}
			}

			c.mu.Unlock()
		}
	}()

	return c
}

// Set adds a new item to the cache with the specified key, value, and
// time-to-live (TTL).
func (c *CacheInstance) Set(key string, value any, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = cacheItem{
		value:  value,
		expiry: time.Now().Add(ttl),
	}
}

// Get retrieves the value associated with the given key from the cache.
func (c *CacheInstance) Get(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if item.isExpired() {
		delete(c.items, key)
		return nil, false
	}

	return item.value, true
}

// Del removes the item with the specified key from the cache.
func (c *CacheInstance) Del(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// Pop removes and returns the item with the specified key from the cache.
func (c *CacheInstance) Pop(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	delete(c.items, key)

	if item.isExpired() {
		return nil, false
	}

	return item.value, true
}

// Len returns the number of items in the cache.
func (c *CacheInstance) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.items)
}

// Clear removes all items from the cache.
func (c *CacheInstance) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]cacheItem)
}
