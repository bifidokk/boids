package main

import (
	"context"
	"fmt"
	"reflect"
)

func main() {
	ch := make(chan interface{})

	go func() {
		for i := 0; i < 3; i++ {
			ch <- i
		}

		close(ch)
	}()

	var res []interface{}

	for v := range onDone(context.Background(), ch) {
		res = append(res, v)
	}

	fmt.Println(res)

	if !reflect.DeepEqual(res, []interface{}{0, 1, 2}) {
		panic("wrong code")
	}
}

func onDone(ctx context.Context, in chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case r, ok := <-in:
				if !ok {
					return
				}

				select { // nested select to be sure that context will not be cancelled during writing
				case <-ctx.Done():
				case out <- r:
				}

			}
		}
	}()

	return out
}
