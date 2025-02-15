package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var m = sync.Map{}

func main() {
	userIds := []string{"user1", "user2", "user3", "user4", "user2", "user1"}
	wg := sync.WaitGroup{}
	wg.Add(len(userIds))

	for _, id := range userIds {
		go func() {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
			result := getOrCompute(id, func() string {
				fmt.Println("Computation in progress")
				return fmt.Sprintf("Result for user %s", id)
			})

			fmt.Println(result)
			wg.Done()
		}()
	}

	wg.Wait()
}

func getOrCompute(key string, compute func() string) string {
	v, ok := m.Load(key)

	if !ok {
		v = compute()
		m.Store(key, v)
	}

	return v.(string)
}
