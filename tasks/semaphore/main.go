package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func download(url string) {
	fmt.Println("Downloading", url)
	time.Sleep(1 * time.Second)
	fmt.Println("Done")
}

const goroutinesLimit = 3

func main() {
	files := []string{"url1", "url2", "url3", "url4", "url5", "url6", "url7", "url8", "url9"}

	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, goroutinesLimit)

	for _, file := range files {
		wg.Add(1)
		fmt.Println(runtime.NumGoroutine())
		semaphore <- struct{}{}

		go func() {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			download(file)
		}()
	}

	wg.Wait()
}
