package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	money          = 100
	lock           = sync.Mutex{}
	moneyDeposited = sync.NewCond(&lock)
)

func main() {
	go income()
	go expense()
	time.Sleep(3000 * time.Millisecond)
	fmt.Println(money)
}

func income() {
	for i := 0; i < 1000; i++ {
		lock.Lock()
		money += 10
		fmt.Println("Balance after income ", money)
		moneyDeposited.Signal()
		lock.Unlock()
		time.Sleep(1 * time.Millisecond)
	}

	fmt.Println("Income done")
}

func expense() {
	for i := 0; i < 1000; i++ {
		lock.Lock()
		for money-20 < 0 {
			moneyDeposited.Wait()
		}

		money -= 20
		fmt.Println("Balance after expense ", money)
		lock.Unlock()
		time.Sleep(1 * time.Millisecond)
	}

	fmt.Println("Expense done")
}
