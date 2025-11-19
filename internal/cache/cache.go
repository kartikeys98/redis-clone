package cache

import "sync"

type Cache struct {
	data map[string]*CacheEntry
	lock sync.RWMutex
	maxSize int
	lruList *LRUList
}

type CacheEntry struct {
	Value string
	lruNode *Node
}

func New(maxSize int) *Cache {
	if maxSize < 0 {
        panic("cache: maxSize cannot be negative")
    }
	cache := &Cache{
		data: make(map[string]*CacheEntry),
		lock: sync.RWMutex{},
		maxSize: maxSize,
		lruList: &LRUList{},
	}
	return cache
}

func (c *Cache) Get(key string) (string, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	entry, ok := c.data[key]
	if !ok {
		return "", false
	}
	c.lruList.MoveToFront(entry.lruNode)
	return entry.Value, true
}

func (c *Cache) Set(key string, value string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	// If the key already exists, update the value and move the node to the front
	if entry, ok := c.data[key]; ok {
		entry.Value = value
		c.lruList.MoveToFront(entry.lruNode)
		return
	}
	// If the cache is full, remove the least recently used node
	if c.maxSize > 0 && len(c.data) >= c.maxSize {
		node := c.lruList.RemoveLRU()
		if node == nil {
			panic("Failed to remove least recently used node")
		}
		delete(c.data, node.Key)
	}
	c.data[key] = &CacheEntry{
		Value: value,
		lruNode: c.lruList.AddToFront(key),
	}
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

// Keys returns all keys in the cache
func (c *Cache) Keys() []string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	keys := make([]string, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
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
