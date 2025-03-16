package main

import (
	"context"
	"fmt"
	"reflect"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	i := 0
	inc := func() interface{} {
		i++
		return i
	}

	out1, out2 := tee(ctx, take(ctx, repeatFn(ctx, inc), 3))

	var res1, res2 []interface{}

	for val1 := range out1 {
		res1 = append(res1, val1)
		res2 = append(res2, <-out2)
	}

	exp := []interface{}{1, 2, 3}

	fmt.Println(res1)
	fmt.Println(res2)

	if !reflect.DeepEqual(res1, exp) || !reflect.DeepEqual(res2, exp) {
		panic("wrong code")
	}
}

func tee(ctx context.Context, in <-chan interface{}) (_, _ <-chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {
		for v := range onDone(ctx, in) {
			select {
			case <-ctx.Done():
				return
			case out1 <- v:
			}

			select {
			case <-ctx.Done():
				return
			case out2 <- v:
			}

			// this is better possible solution to write to two channels without blocking.
			// at first local out1 and out2 point to their respective channels, after writing to channel we change local variable to nil
			// and in this case select skip it. So we write into both channels without blocking
			//var out1, out2 = out1, out2
			//for i := 0; i < 2; i++ {
			//	select {
			//	case <-ctx.Done():
			//	case out1 <- v:
			//		out1 = nil
			//	case out2 <- v:
			//		out2 = nil
			//	}
			//}
		}

		close(out1)
		close(out2)
	}()

	return out1, out2
}

func onDone(ctx context.Context, in <-chan interface{}) <-chan interface{} {
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

func repeatFn(ctx context.Context, fn func() interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case out <- fn():
			}
		}
	}()
	return out
}

func take(ctx context.Context, in <-chan interface{}, num int) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for i := 0; i < num; i++ {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- v
			}
		}
	}()
	return out
}
