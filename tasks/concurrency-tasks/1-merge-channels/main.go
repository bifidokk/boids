package main

import (
	"fmt"
	"sync"
)

// merge all channels into one
func merge(cs ...<-chan int) <-chan int {
	out := make(chan int)
	wg := sync.WaitGroup{}

	wg.Add(len(cs))

	for _, c := range cs {
		go func(c <-chan int) {
			for n := range c {
				out <- n
			}

			wg.Done()
		}(c)

	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// fill channel with numbers 0, n-1
func fillChan(n int) <-chan int {
	c := make(chan int)

	go func() {
		for i := 0; i < n; i++ {
			c <- i
		}

		close(c)
	}()

	return c
}

func main() {
	a := fillChan(2) // 0, 1
	b := fillChan(3) // 0, 1, 2
	c := fillChan(4) // 0, 1, 2, 3

	d := merge(a, b, c)

	for n := range d {
		fmt.Println(n)
	}
}
