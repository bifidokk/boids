package main

import "sync"

func main() {
	channels := make([]chan int, 10)
	for i := range channels {
		channels[i] = make(chan int)
	}

	for i := range channels {
		go func(i int) {
			channels[i] <- i
			close(channels[i])
		}(i)
	}

	for v := range merge(channels...) {
		println(v)
	}

}

func merge(channels ...chan int) chan int {
	res := make(chan int)

	merger := func(ch chan int) {
		for {
			select {
			case v, ok := <-ch:
				if !ok {
					return
				}

				res <- v
			}
		}

	}

	wg := sync.WaitGroup{}
	for _, ch := range channels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			merger(ch)
		}()
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	return res
}
