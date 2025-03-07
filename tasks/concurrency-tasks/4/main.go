package main

import (
	"context"
	"math/rand"
)

// The repeatFn function calls the fn function infinitely and writes its result to the returned channel.
// Terminates early if the context is canceled.
// The take function reads at most num from the in channel while in is open and writes the value to the returned channel.
// Terminates early if the context is canceled.

func repeatFn(ctx context.Context, fn func() interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case out <- fn():
			}
		}
	}()

	return out

}

func take(ctx context.Context, in <-chan interface{}, num int) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for i := 0; i < num; i++ {
			select {
			case <-ctx.Done():
				return
			case r, ok := <-in:
				if !ok {
					return
				}

				out <- r
			}
		}
	}()

	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	randFn := func() interface{} {
		return rand.Int()
	}

	var result []interface{}

	for num := range take(ctx, repeatFn(ctx, randFn), 3) {
		result = append(result, num)
	}

	if len(result) != 3 {
		panic("wrong")
	}
}
