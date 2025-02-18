package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

var counter atomic.Int64

func SimulateRequest(ctx context.Context) (int64, error) {
	start := time.Now()

	defer func() {
		fmt.Printf("Simulated request took %v\n", time.Since(start))
	}()

	requestResultChannel := make(chan int64)

	go func() {
		time.Sleep(time.Duration(rand.Int63n(5)) * time.Second)
		counter.Add(1)
		requestResultChannel <- counter.Load()
	}()

	select {
	case <-ctx.Done():
		{
			fmt.Println("Cancelled")
			return 0, ctx.Err()
		}
	case result := <-requestResultChannel:
		return result, nil
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	val, err := SimulateRequest(ctx)

	fmt.Println(val, err)
}
