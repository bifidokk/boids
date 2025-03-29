package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const (
	produceCount = 3
	produceStop  = 10
)

func main() {
	pipe := make(chan int)

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	for i := 0; i < produceCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			produce(ctx, pipe)
		}()
	}

	for r := range pipe {
		fmt.Println(r)
		if r == produceStop {
			cancel()
			break
		}
	}

	wg.Wait()
	close(pipe)

	fmt.Println("main finished")
}

func produce(ctx context.Context, pipe chan<- int) {
	for i := 0; ; i++ {
		select {
		case <-ctx.Done():
			time.Sleep(3 * time.Second)
			fmt.Println("produce finished")
			return
		case pipe <- i:
		}
	}
}
