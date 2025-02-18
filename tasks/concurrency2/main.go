package main

import (
	"context"
	"fmt"
	"time"
)

func getDiscount() float64 {
	time.Sleep(5 * time.Second)

	return 12.0
}

func getDiscountWithTimeout(ctx context.Context) (float64, error) {
	resultChan := make(chan float64)

	go func() {
		discount := getDiscount()
		resultChan <- discount
		close(resultChan)
	}()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	case result := <-resultChan:
		return result, nil
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	discount, err := getDiscountWithTimeout(ctx)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Discount is", discount)
}
