package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

var (
	money int32 = 100
)

func main() {
	go income()
	go expense()
	time.Sleep(3000 * time.Millisecond)
	fmt.Println(money)
}

func income() {
	for i := 0; i < 1000; i++ {
		atomic.AddInt32(&money, 10)
		time.Sleep(1 * time.Millisecond)
	}

	fmt.Println("Income done")
}

func expense() {
	for i := 0; i < 1000; i++ {
		atomic.AddInt32(&money, -10)
		time.Sleep(1 * time.Millisecond)
	}

	fmt.Println("Expense done")
}
