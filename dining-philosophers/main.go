package main

import (
	"fmt"
	"sync"
	"time"
)

type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

var philosophers = [5]Philosopher{
	{name: "P1", rightFork: 4, leftFork: 0},
	{name: "P2", rightFork: 0, leftFork: 1},
	{name: "P3", rightFork: 1, leftFork: 2},
	{name: "P4", rightFork: 2, leftFork: 3},
	{name: "P5", rightFork: 3, leftFork: 4},
}

var (
	hunger    = 3
	eatTime   = 0 * time.Second
	idleTime  = 0 * time.Second
	sleepTime = 0 * time.Second

	philosophersOrderFinishedMutex sync.Mutex
	philosophersOrderFinished      []string
)

func main() {
	start := time.Now()
	fmt.Println("Starting...")

	dine()

	elapsed := time.Since(start)
	fmt.Println("Elapsed: ", elapsed)
}

func dine() {
	wg := sync.WaitGroup{}
	wg.Add(len(philosophers))

	seated := sync.WaitGroup{}
	seated.Add(len(philosophers))

	var forks = make(map[int]*sync.Mutex)

	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	for i := 0; i < len(philosophers); i++ {
		go processDining(philosophers[i], &wg, forks, &seated)
	}

	wg.Wait()
	fmt.Println("Order ", philosophersOrderFinished)
}

func processDining(philosopher Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seated *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("%s is seated st the table\n", philosopher.name)
	seated.Done()
	seated.Wait()

	for i := hunger; i > 0; i-- {
		if philosopher.leftFork > philosopher.rightFork {
			forks[philosopher.rightFork].Lock()
			fmt.Printf("%s takes the right fork %d\n", philosopher.name, philosopher.rightFork)
			forks[philosopher.leftFork].Lock()
			fmt.Printf("%s takes the left fork %d\n", philosopher.name, philosopher.leftFork)
		} else {
			forks[philosopher.leftFork].Lock()
			fmt.Printf("%s takes the left fork %d\n", philosopher.name, philosopher.leftFork)
			forks[philosopher.rightFork].Lock()
			fmt.Printf("%s takes the right fork %d\n", philosopher.name, philosopher.rightFork)
		}

		fmt.Printf("%s is eating\n", philosopher.name)
		time.Sleep(eatTime)

		fmt.Printf("%s is thinking\n", philosopher.name)
		time.Sleep(idleTime)

		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()

		fmt.Printf("%s puts down forks\n", philosopher.name)
	}

	philosophersOrderFinishedMutex.Lock()
	philosophersOrderFinished = append(philosophersOrderFinished, philosopher.name)
	philosophersOrderFinishedMutex.Unlock()

	fmt.Printf("%s is finished\n", philosopher.name)
}
