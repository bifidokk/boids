package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

const timeout = 100 * time.Millisecond

func main() {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	err := executeTaskWithTimeout(ctx)

	if err != nil {
		panic(err)
	}

	fmt.Println("task done")
}

func executeTaskWithTimeout(ctx context.Context) error {
	c := make(chan struct{})

	go func() {
		executeTask()
		c <- struct{}{}
		close(c)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c:
		return nil
	}
}

func executeTask() {
	time.Sleep(time.Duration(rand.Intn(3)) * timeout)
}
