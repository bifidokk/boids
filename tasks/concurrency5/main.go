package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type ParkingLot struct {
	slots chan struct{}
}

func (p *ParkingLot) Park(carID int64) {
	p.slots <- struct{}{}

	fmt.Printf("parking started %d\n", carID)
	time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
	fmt.Printf("parking finished %d\n", carID)

	<-p.slots
}

func main() {
	parkingLot := &ParkingLot{
		slots: make(chan struct{}, 3),
	}

	wg := sync.WaitGroup{}

	carIDs := []int64{1, 2, 3, 4, 5, 6}

	for _, id := range carIDs {
		wg.Add(1)

		go func(carID int64) {
			defer wg.Done()
			parkingLot.Park(carID)
		}(id)
	}

	wg.Wait()
	close(parkingLot.slots)
}
