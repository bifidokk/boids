package main

import (
	"fmt"
	"time"
)

func UpdateProductStock() <-chan map[string]int {
	updates := make(chan map[string]int)

	go func() {
		defer close(updates)

		currentStock := map[string]int{
			"apples":  50,
			"oranges": 40,
			"bananas": 10,
			"grapes":  15,
		}

		for i := 0; i < 5; i++ {
			// make a copy. Otherwise, we send a pointer, so we will have the same final result in all history records
			newStock := make(map[string]int)

			for product, quantity := range currentStock {
				newStock[product] = int(float64(quantity) * 0.95)
			}

			updates <- newStock
			currentStock = newStock

			time.Sleep(150 * time.Millisecond)
		}
	}()

	return updates
}

func main() {
	stream := UpdateProductStock()

	var stockHistory []map[string]int

	for stock := range stream {
		stockHistory = append(stockHistory, stock)
	}

	for i, stock := range stockHistory {
		fmt.Printf("Iteration %d: %v\n", i+1, stock)
	}
}
