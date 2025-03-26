package main

import (
	"errors"
	"sync"
)

// Write a function Run that concurrently executes functions fs and waits for their completion.
// If one or more functions from fs complete with an error, Run returns any of these errors.

type fn func() error

func main() {
	expErr := errors.New("error")

	funcs := []fn{
		func() error { return nil },
		func() error { return nil },
		func() error { return expErr },
		func() error { return nil },
	}

	if err := Run(funcs...); !errors.Is(err, expErr) {
		panic("wrong code")
	}
}

func Run(fs ...fn) error {
	ch := make(chan error)
	wg := sync.WaitGroup{}

	for _, f := range fs {
		wg.Add(1)

		go func() {
			defer wg.Done()

			result := f()
			ch <- result
		}()
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	var firstErr error
	for err := range ch {
		if firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}
