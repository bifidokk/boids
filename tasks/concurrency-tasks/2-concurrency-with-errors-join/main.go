package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Result struct {
	msg string
	err error
}

func download(urls []string) ([]string, error) {
	wg := sync.WaitGroup{}
	ch := make(chan Result, len(urls))

	for _, url := range urls {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()
			result := fakeDownload(url)

			ch <- result
		}(url)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	results := []string{}
	var err error

	for result := range ch {
		if result.err != nil {
			err = errors.Join(err, result.err)
		}

		results = append(results, result.msg)
	}

	return results, err
}

func fakeDownload(url string) Result {
	fmt.Println("Downloading", url)
	time.Sleep(1 * time.Second)

	resp, err := http.Get(url)
	if err != nil {
		return Result{
			err: err,
		}
	}

	defer resp.Body.Close()

	return Result{
		msg: "Result for " + url + " (status code " + strconv.Itoa(resp.StatusCode) + ")",
	}
}

func main() {
	urls := []string{
		"https://www.google.com",
		"https://www.facebook.com",
		"https://x.com/home",
		"https://www.ramonescore.com",
		"https://www.ramonescore2.com",
		"https://www.youtube.com",
		"https://www.reddit.com",
	}

	results, err := download(urls)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Results", results)
}
