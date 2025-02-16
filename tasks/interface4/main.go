package main

import (
	"fmt"
	"sync"
)

type Cache struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	c.data[key] = value
	c.mu.Unlock()
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	v, ok := c.data[key]
	c.mu.RUnlock()
	return v, ok
}

func main() {
	cache := &Cache{data: make(map[string]interface{})}

	cache.Set("foo", "bar")
	cache.Set("bar", "baz")

	fmt.Println(cache.Get("foo"))
	fmt.Println(cache.Get("bar"))
	fmt.Println(cache.Get("baz"))
}
