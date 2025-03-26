package main

func main() {
	first := make(chan int)
	last := make(<-chan int)

	n := 10

	last = inc(first) // run 10 times to increase 0 to 10
	for i := 1; i < n; i++ {
		last = inc(last)
	}

	first <- 0
	close(first)
	if n != <-last {
		panic("wrong code")
	}
}

func inc(in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for i := range in {
			out <- i + 1
		}
	}()

	return out
}
