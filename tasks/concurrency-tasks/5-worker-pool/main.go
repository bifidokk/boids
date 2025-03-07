package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func worker(f func(int) int, jobs <-chan int, results chan<- int) {
	for x := range jobs {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		fmt.Println("Write ", x)
		results <- f(x)
	}
}

const numJobs = 5
const numWorkers = 3

func main() {
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	multiplier := func(x int) int {
		return x * 10
	}

	wg := sync.WaitGroup{}
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func() {
			defer wg.Done()
			worker(multiplier, jobs, results)
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		for i := 0; i < numJobs; i++ {
			jobs <- i
		}

		close(jobs)
	}()

	done := make(chan struct{})

	go func() {
		for res := range results {
			fmt.Println(res)
		}

		close(done)
	}()

	<-done

	fmt.Println("DONE")
}
