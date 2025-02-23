package main

import (
	"fmt"
	"hash/fnv"
	"sync"
)

type Cache interface {
	Set(key string, value string)
	Get(key string) (string, bool)
}

type Shard struct {
	data map[string]string
	mu   sync.RWMutex
}

type MemoryCache struct {
	shards []*Shard
}

func NewMemoryCache(size int64) *MemoryCache {
	shards := make([]*Shard, size)

	for i := range shards {
		shards[i] = &Shard{data: make(map[string]string)}
	}

	return &MemoryCache{
		shards: shards,
	}
}

func (c *MemoryCache) Set(key string, value string) {
	shard := c.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.data[key] = value
}

func (c *MemoryCache) Get(key string) (string, bool) {
	shard := c.getShard(key)

	shard.mu.RLock()
	defer shard.mu.RUnlock()

	v, ok := shard.data[key]

	return v, ok
}

func (c *MemoryCache) getShard(key string) *Shard {
	h := fnv.New32a()
	_, err := h.Write([]byte(key))

	if err != nil {
		return nil
	}

	return c.shards[int(h.Sum32())%len(c.shards)]
}

// it's possible to use consistent hashing to avoid problems if we add more shards
func main() {
	cache := NewMemoryCache(100)
	cache.Set("foo", "bar")

	val, ok := cache.Get("foo")
	if !ok {
		fmt.Println("foo not found")
		return
	}

	fmt.Println("foo:", val)
}
