package cache

import (
	"sync"
	"time"
)

type Cache struct {
	data        map[string]*CacheEntry
	lock        sync.RWMutex
	maxSize     int
	lruList     *LRUList
	stopCleanup chan struct{} // Channel to signal background cleanup goroutine to stop
}

type CacheEntry struct {
	Value      string
	lruNode    *Node
	ExpiryTime time.Time
}

func New(maxSize int) *Cache {
	if maxSize < 0 {
		panic("cache: maxSize cannot be negative")
	}
	cache := &Cache{
		data:    make(map[string]*CacheEntry),
		lock:    sync.RWMutex{},
		maxSize: maxSize,
		lruList: &LRUList{},
		stopCleanup: make(chan struct{}),
	}
	go cache.backgroundCleanup()
	return cache
}

func (c *Cache) backgroundCleanup() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.cleanupExpiredKeys()
		case <-c.stopCleanup:
			return
		}
	}
}

func (c *Cache) cleanupExpiredKeys() {
	c.lock.Lock()
	defer c.lock.Unlock()
	// Check if any TTL is expired - if so, delete the key
	for key, entry := range c.data {
		if entry.ExpiryTime.Before(time.Now()) && !entry.ExpiryTime.IsZero() {
			c.deleteWithoutLocking(key)
		}
	}
}

func (c *Cache) Close() {
	close(c.stopCleanup)
}

func (c *Cache) Get(key string) (string, bool) {
	c.lock.Lock() // Why not RLock? Because we need to update the LRU list and delete the key if it's expired
	defer c.lock.Unlock()
	entry, ok := c.data[key]
	if !ok {
		return "", false
	}
	if entry.ExpiryTime.IsZero() || entry.ExpiryTime.After(time.Now()) {
		c.lruList.MoveToFront(entry.lruNode)
		return entry.Value, true
	}
	c.deleteWithoutLocking(key)
	return "", false
}

func (c *Cache) SetWithTTL(key string, value string, ttl time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()
	// If the key already exists, update the value and move the node to the front
	if entry, ok := c.data[key]; ok {
		entry.Value = value
		entry.ExpiryTime = time.Time{}
		if ttl > 0 {
			entry.ExpiryTime = time.Now().Add(ttl)
		}
		c.lruList.MoveToFront(entry.lruNode)
		return
	}
	// If the cache is full, remove the least recently used node
	if c.maxSize > 0 && len(c.data) >= c.maxSize {
		// Check if any TTL is expired - if so, delete the key
		for key, entry := range c.data {
			if entry.ExpiryTime.Before(time.Now()) && !entry.ExpiryTime.IsZero() {
				c.deleteWithoutLocking(key)
			}
		}
		// If the cache is still full, remove the least recently used node
		if len(c.data) >= c.maxSize {
			node := c.lruList.RemoveLRU()
			if node == nil {
				panic("Failed to remove least recently used node")
			}
			delete(c.data, node.Key)
		}
	}
	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}
	c.data[key] = &CacheEntry{
		Value:      value,
		lruNode:    c.lruList.AddToFront(key),
		ExpiryTime: expiresAt,
	}
}

func (c *Cache) Set(key string, value string) {
	c.SetWithTTL(key, value, 0)
}

func (c *Cache) Delete(key string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	entry, ok := c.data[key]
	if !ok {
		return false
	}
	if entry.lruNode == nil {
		panic("Failed to remove node from LRU list")
	}
	c.lruList.Remove(entry.lruNode)
	delete(c.data, key)
	return true
}

func (c *Cache) deleteWithoutLocking(key string) {
	entry, ok := c.data[key]
	if !ok {
		return
	}
	if entry.lruNode == nil {
		panic("Failed to remove node from LRU list")
	}
	c.lruList.Remove(entry.lruNode)
	delete(c.data, key)
}

// Keys returns all keys in the cache
func (c *Cache) Keys() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	keys := make([]string, 0, len(c.data))
	for key, entry := range c.data {
		if entry.ExpiryTime.IsZero() || entry.ExpiryTime.After(time.Now()) {
			keys = append(keys, key)
		}
	}
	return keys
}

// Flush removes all keys from the cache
func (c *Cache) Flush() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data = make(map[string]*CacheEntry)
	c.lruList = &LRUList{
		Head: nil,
		Tail: nil,
		Size: 0,
	}
}

// Size returns the number of keys in the cache
func (c *Cache) Size() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.data)
}
