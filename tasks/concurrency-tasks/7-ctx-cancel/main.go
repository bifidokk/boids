package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type result struct {
	msg string
	err error
}

type search func() *result

type replicas []search

func fakeSearch(kind string) search {
	return func() *result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		return &result{
			msg: fmt.Sprintf("%q result", kind),
		}
	}
}

func getFirstResult(ctx context.Context, replicas replicas) *result {
	c := make(chan *result)

	// create a new context with cancel to cancel all goroutines when the first result is received
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // cancel will be called after the first result is received

	for _, replica := range replicas {
		go func(replica search) {
			select {
			case <-ctx.Done(): // ctx is cancelled, that means goroutine is finished without calling the case c <- replica()
			case c <- replica():
			}

		}(replica)
	}

	select {
	case <-ctx.Done():
		return &result{err: ctx.Err()}
	case r := <-c:
		return r
	}
}

func getResults(ctx context.Context, replicaKinds []replicas) []*result {
	c := make(chan *result)
	wg := sync.WaitGroup{}

	for _, rep := range replicaKinds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c <- getFirstResult(ctx, rep)
		}()
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	var resultSlice []*result

	for r := range c {
		resultSlice = append(resultSlice, r)
	}

	return resultSlice
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 50*time.Millisecond)

	replicaKinds := []replicas{
		replicas{fakeSearch("web1"), fakeSearch("web2")},
		replicas{fakeSearch("image1"), fakeSearch("image2")},
		replicas{fakeSearch("video1"), fakeSearch("video2")},
	}

	for _, res := range getResults(ctx, replicaKinds) {
		fmt.Println(res.msg, res.err)
	}
}
