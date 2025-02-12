package main

import (
	"fmt"
	"sync"
)

type ConcurrentMap struct {
	m  map[string]string
	mu sync.RWMutex
}

// read write locks divided for better performance if there are a lot of reads
func (m *ConcurrentMap) GetOrCreate(key string, value string) string {
	m.mu.RLock()
	val, ok := m.m[key]
	m.mu.RUnlock()

	if ok {
		return val
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if val, ok = m.m[key]; ok {
		return val
	}

	m.m[key] = value

	return value
}

func main() {
	m := NewConcurrentMap()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		val := m.GetOrCreate("key1", "value1")
		fmt.Println("Goroutine 1 ", val)
	}()

	go func() {
		defer wg.Done()

		val := m.GetOrCreate("key1", "value2")
		fmt.Println("Goroutine 2 ", val)
	}()

	wg.Wait()
}

func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		m:  make(map[string]string),
		mu: sync.RWMutex{},
	}
}
