package main

import "fmt"

func mergeSorted(a, b <-chan int) chan int {
	out := make(chan int)

	go func() {
		for {
			ra, ok1 := <-a
			rb, ok2 := <-b

			if !ok1 && !ok2 {
				break
			}

			if ra < rb {
				if ok1 {
					out <- ra
				}

				if ok2 {
					out <- rb
				}
			} else {
				if ok2 {
					out <- rb
				}

				if ok1 {
					out <- ra
				}
			}
		}

		close(out)
	}()

	return out
}

func fillChanA(c chan int) {
	c <- 1
	c <- 2
	c <- 4
	close(c)
}

func fillChanB(c chan int) {
	c <- -1
	c <- 4
	c <- 5
	c <- 8
	close(c)
}

func main() {
	a := make(chan int)
	b := make(chan int)

	go fillChanA(a)
	go fillChanB(b)

	c := mergeSorted(a, b)

	for r := range c {
		fmt.Println(r)
	}
}
