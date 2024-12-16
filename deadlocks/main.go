package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	lock1 = sync.Mutex{}
	lock2 = sync.Mutex{}
)

func blueRobot() {
	for {
		fmt.Println("blueRobot: acquiring lock1")
		lock1.Lock()
		fmt.Println("blueRobot: acquiring lock2")
		lock2.Lock()
		fmt.Println("blueRobot: both locks acquired")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("blueRobot: locks released")
	}
}

func redRobot() {
	for {
		fmt.Println("redRobot: acquiring lock2")
		lock2.Lock()
		fmt.Println("redRobot: acquiring lock1")
		lock1.Lock()
		fmt.Println("redRobot: both locks acquired")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("redRobot: locks released")
	}
}

func main() {
	go redRobot()
	go blueRobot()
	time.Sleep(time.Second * 20)
	fmt.Println("done")
}
