package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

//The waiter interface should:
//
//Concurrently run functions passed to run with specified context
//Set parallel execution limit via maxParallel parameter when creating waiter through newGroupWait
//Return an error from wait if any function in run returns an error
//Return a combined error from multiple failed tasks using errors.Join

type waiter interface {
	wait() error
	run(ctx context.Context, f func(ctx context.Context) error)
}

type waitGroup struct {
	semaphore   chan struct{}
	maxParallel int
	wg          sync.WaitGroup
	mu          sync.Mutex
	errs        []error
}

func (g *waitGroup) wait() error {
	g.wg.Wait()

	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.errs) == 0 {
		return nil
	}

	return errors.Join(g.errs...)
}

func (g *waitGroup) run(ctx context.Context, f func(ctx context.Context) error) {
	g.semaphore <- struct{}{} // acquire the slot
	g.wg.Add(1)

	go func() {
		defer func() {
			<-g.semaphore // release the slot
			g.wg.Done()
		}()

		if err := f(ctx); err != nil {
			g.mu.Lock()
			g.errs = append(g.errs, err)
			g.mu.Unlock()
		}
	}()

}

func newGroupWait(maxParallel int) waiter {
	return &waitGroup{
		maxParallel: maxParallel,
		semaphore:   make(chan struct{}, maxParallel),
		errs:        []error{},
	}
}

func main() {
	g := newGroupWait(2)

	ctx := context.Background()
	expErr1 := errors.New("error 1")
	expErr2 := errors.New("error 2")

	g.run(ctx, func(ctx context.Context) error {
		return nil
	})

	g.run(ctx, func(ctx context.Context) error {
		return expErr2
	})

	g.run(ctx, func(ctx context.Context) error {
		return expErr1
	})

	err := g.wait()

	fmt.Printf("%v", err)

	if !errors.Is(err, expErr1) || !errors.Is(err, expErr2) {
		panic("wrong code")
	}
}
