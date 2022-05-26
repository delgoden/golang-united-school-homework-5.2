package cache

import (
	"sync"
	"time"
)

type Cache struct {
	sync.RWMutex
	data map[string]string
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]string)}
}

func (c *Cache) Get(key string) (string, bool) {
	c.RLock()
	if value, ok := c.data[key]; ok {
		return value, ok
	}
	c.RUnlock()
	return "", false
}

func (c *Cache) Put(key, value string) {
	c.Lock()
	c.data[key] = value
	c.Unlock()
}

func (c *Cache) Keys() []string {
	keys := make([]string, 0, len(c.data))
	c.RLock()
	for key := range c.data {
		keys = append(keys, key)
	}
	c.RUnlock()
	return keys
}

func (c *Cache) PutTill(key, value string, deadline time.Time) {
	c.Lock()
	c.data[key] = value
	c.Unlock()
	go func(deadline time.Time, key string) {
		time.Sleep(deadline.Sub(time.Now()))
		c.Lock()
		delete(c.data, key)
		c.Unlock()
	}(deadline, key)
}
