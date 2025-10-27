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
