package main

import (
	"fmt"
	"time"
)

// in this case there will be deadlock because single goroutine will write to the channel and select will write there after tick
// to avoid it we can create buffered channel
func main() {
	ch := make(chan bool)

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Single goroutine is done")
		ch <- false
	}()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Tick")
			ch <- true
		case value := <-ch:
			fmt.Println("Channel is done", value)
			return
		}
	}
}
