package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const numberOfOrders = 10

var successAmount, failedAmount, total int

type Producer struct {
	data chan Order
	done chan bool
}

type Order struct {
	orderNumber int
	message     string
	success     bool
}

func main() {
	fmt.Println("Starting producer")

	ordersJob := &Producer{make(chan Order), make(chan bool)}

	var wg sync.WaitGroup
	wg.Add(1)

	go processOrders(ordersJob, &wg)

	for i := range ordersJob.data {
		if i.orderNumber <= numberOfOrders {
			if i.success {
				fmt.Printf("Order #%d has been successfully received by customer\n", i.orderNumber)
			} else {
				fmt.Printf("Order #%d has not been successfully received :(\n", i.orderNumber)
			}
		} else {
			ordersJob.Close()
		}
	}

	wg.Wait()
}

func processOrders(producer *Producer, wg *sync.WaitGroup) {
	defer wg.Done()

	i := 0

	for {
		order := createOrder(i)

		if order != nil {
			i = order.orderNumber

			select {
			case producer.data <- *order:
			case <-producer.done:
				close(producer.data)
				close(producer.done)
				return
			}
		}
	}
}

func createOrder(i int) *Order {
	i++

	if i > numberOfOrders {
		return nil
	}

	delay := rand.Intn(5) + 1
	rnd := rand.Intn(12) + 1
	fmt.Println("Created order", i)
	success := true
	message := fmt.Sprintf("Order #%d has been successfully processed", i)

	if rnd <= 4 {
		failedAmount++
		success = false
		message = fmt.Sprintf("Order #%d has been failed", i)
	} else {
		successAmount++
	}

	total++

	time.Sleep(time.Duration(delay) * time.Second)

	return &Order{i, message, success}
}

func (p *Producer) Close() {
	p.done <- true
}
