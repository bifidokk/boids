package main

import (
	"fmt"
	"math/rand"
)

type Transaction struct {
	ID     int64
	Amount float64
}

func create(count int) <-chan Transaction {
	output := make(chan Transaction)

	go func() {
		for i := 0; i < count; i++ {
			output <- Transaction{
				ID:     int64(i),
				Amount: rand.Float64()*200 - 100,
			}
		}

		close(output)
	}()

	return output
}

func filter(input <-chan Transaction) <-chan Transaction {
	output := make(chan Transaction)

	go func() {
		for t := range input {
			if t.Amount >= 0 {
				output <- t
			}
		}

		close(output)
	}()

	return output
}

func convert(input <-chan Transaction) <-chan Transaction {
	output := make(chan Transaction)

	go func() {
		for t := range input {
			t.Amount *= 0.8
			output <- t
		}

		close(output)
	}()

	return output
}

func store(input <-chan Transaction) {
	for t := range input {
		fmt.Println("store", t.ID, t.Amount)
	}
}

func main() {
	transactions := create(1000)
	filtered := filter(transactions)
	converted := convert(filtered)
	store(converted)
}
