package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func streamNumber(done <-chan bool, generateNumber func() int) <-chan int {
	c := make(chan int)

	go func() {
		defer close(c)

		for {
			select {
			case <-done:
				return
			case c <- generateNumber():
			}
		}
	}()

	return c
}

func primeFinder(done <-chan bool, stream <-chan int) <-chan int {
	isPrime := func(num int) bool {
		for i := num - 1; i > 1; i-- {
			if num%i == 0 {
				return false
			}
		}

		return true
	}

	primes := make(chan int)

	go func() {
		defer close(primes)

		for {
			select {
			case <-done:
				return
			case prime := <-stream:
				if isPrime(prime) {
					primes <- prime
				}
			}
		}
	}()

	return primes
}

func fanIn(done <-chan bool, primeChannels []<-chan int) <-chan int {
	wg := sync.WaitGroup{}

	fanInStream := make(chan int)
	transfer := func(c <-chan int) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case fanInStream <- i:
			}
		}
	}

	for _, primeChannel := range primeChannels {
		wg.Add(1)
		go transfer(primeChannel)
	}

	go func() {
		wg.Wait()
		close(fanInStream)
	}()

	return fanInStream
}

func take(done <-chan bool, stream <-chan int, n int) <-chan int {
	c := make(chan int)

	go func() {
		defer close(c)

		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case c <- <-stream:
			}
		}
	}()

	return c
}

func main() {
	start := time.Now()

	done := make(chan bool)
	defer close(done)

	randIntStream := streamNumber(done, generateNumber)

	cpuCount := runtime.NumCPU()
	primeFinderChannels := make([]<-chan int, cpuCount)
	fmt.Println("fan out to", cpuCount, "CPU")

	// fan out
	for i := 0; i < cpuCount; i++ {
		primeFinderChannels[i] = primeFinder(done, randIntStream)
	}

	// fan in
	fanInStream := fanIn(done, primeFinderChannels)

	for i := range take(done, fanInStream, 10) {
		fmt.Println(i)
	}

	elapsed := time.Since(start)
	fmt.Printf("elapsed: %s\n", elapsed)
}

func generateNumber() int {
	randomNumber := rand.Intn(1000000000)

	return randomNumber
}
