package cache

import "sync"

type Cache struct {
	data map[string]string
	lock sync.RWMutex
}

func New() *Cache {
	cache := &Cache{
		data: make(map[string]string),
		lock: sync.RWMutex{},
	}
	return cache
}

func (c *Cache) Get(key string) (string, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

func (c *Cache) Set(key string, value string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.data[key] = value
}

func (c *Cache) Delete(key string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.data[key]
	delete(c.data, key)
	return ok
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
	c.data = make(map[string]string)
}

// Size returns the number of keys in the cache
func (c *Cache) Size() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return len(c.data)
}
