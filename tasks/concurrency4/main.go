package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func LongCalculation(n int) int {
	secondsToSleep := rand.Float64() * float64(n)
	fmt.Println(secondsToSleep, " seconds to sleep")
	time.Sleep(time.Duration(secondsToSleep) * time.Second)

	return n + 1
}

// we can also use a struct
var cache = map[int]int{}
var mutex = &sync.RWMutex{}

func CachedLongCalculation(n int) int {
	fmt.Println("Calculate for ", n)
	mutex.RLock()
	found, ok := cache[n] // write lock for read operations
	mutex.RUnlock()

	if !ok {
		value := LongCalculation(n)
		mutex.Lock()
		cache[n] = value
		mutex.Unlock()

		return value
	}

	// mutex.Unlock() redundant
	return found
}

func main() {
	nums := []int{22, 18, 10, 22}
	wg := sync.WaitGroup{}

	for _, n := range nums {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(rand.Intn(20)) * time.Second)
			value := CachedLongCalculation(n)
			fmt.Printf("%d -> %d\n", n, value)
			wg.Done()
		}()
	}

	wg.Wait()
}
