package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

func main() {
	wg := sync.WaitGroup{}
	var frequency [26]int32

	start := time.Now()

	for i := 1000; i < 1200; i++ {
		wg.Add(1)
		go countLetters(fmt.Sprintf("https://www.rfc-editor.org/rfc/rfc%d.txt", i), &frequency, &wg)
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Println("elapsed:", elapsed)

	for i, f := range frequency {
		fmt.Printf("%s -> %d\n", string(allLetters[i]), f)
	}
}

func countLetters(url string, frequency *[26]int32, wg *sync.WaitGroup) {
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	for i := 0; i < 20; i++ {
		for _, b := range body {
			c := strings.ToLower(string(b))
			index := strings.Index(allLetters, c)

			if index != -1 {
				atomic.AddInt32(&frequency[index], 1)
			}
		}
	}

	wg.Done()
}
