package main

import (
	"context"
	"fmt"
)

func generator(ctx context.Context, in ...int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for _, i := range in {
			select {
			case <-ctx.Done():
				return
			case out <- i:
			}
		}
	}()

	return out
}

func squarer(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for i := range in {
			select {
			case <-ctx.Done():
				return
			case out <- i * i:
			}
		}
	}()

	return out
}

func main() {
	ctx := context.Background()
	pipeline := squarer(ctx, generator(ctx, 1, 2, 3))

	for x := range pipeline {
		fmt.Println(x)
	}
}
