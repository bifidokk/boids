package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	capacity    = 10
	arrivalRate = 100
	duration    = 1000 * time.Millisecond
	timeOpen    = 10 * time.Second
)

type Barbershop struct {
	capacity        int
	duration        time.Duration
	numberOfBarbers int
	doneChan        chan bool
	clientsChan     chan string
	isOpen          bool
}

func main() {
	fmt.Println("Starting...")

	clientsChan := make(chan string, capacity)
	doneChan := make(chan bool)

	shop := &Barbershop{
		capacity:        capacity,
		duration:        duration,
		numberOfBarbers: 0,
		doneChan:        doneChan,
		clientsChan:     clientsChan,
		isOpen:          true,
	}

	fmt.Println("The shop is open...")

	shop.addBarber("Tony")
	shop.addBarber("Frank")
	shop.addBarber("Dan")

	shopClosingChan := make(chan bool)
	shopClosedChan := make(chan bool)
	go func() {
		<-time.After(timeOpen)
		shopClosingChan <- true
		shop.closeShopForDay()
		shopClosedChan <- true
	}()

	// add clients

	i := 1
	go func() {
		for {
			randomMilliseconds := rand.Int() % (2 * arrivalRate)
			select {
			case <-shopClosingChan:
				return
			case <-time.After(time.Duration(randomMilliseconds) * time.Millisecond):
				shop.addClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()

	<-shopClosedChan
}

func (shop *Barbershop) addBarber(barber string) {
	shop.numberOfBarbers++

	go func() {
		isSleeping := false

		fmt.Printf("%s checks for clients\n", barber)

		for {
			if len(shop.clientsChan) == 0 {
				fmt.Printf("There are no clients. %s sleeps\n", barber)
				isSleeping = true
			}

			client, ok := <-shop.clientsChan

			if ok {
				if isSleeping {
					fmt.Printf("%s wakes %s up\n", client, barber)
					isSleeping = false
				}

				shop.cutHair(barber, client)

			} else {
				shop.sendBarberHome(barber)
				return
			}
		}

	}()
}

func (shop *Barbershop) cutHair(barber string, client string) {
	fmt.Printf("%s cutting %s\n", barber, client)

	time.Sleep(duration)

	fmt.Printf("%s is finished cutting %s\n", barber, client)
}

func (shop *Barbershop) sendBarberHome(barber string) {
	fmt.Printf("%s going home\n", barber)
	shop.doneChan <- true
}

func (shop *Barbershop) closeShopForDay() {
	fmt.Println("Closing shop for day...")

	close(shop.clientsChan)
	shop.isOpen = false

	for a := 1; a <= shop.numberOfBarbers; a++ {
		<-shop.doneChan
	}

	close(shop.doneChan)

	fmt.Println("The shop is closed for day...")
}

func (shop *Barbershop) addClient(name string) {
	fmt.Printf("Adding client %s\n", name)

	if shop.isOpen {
		select {
		case shop.clientsChan <- name:
			fmt.Printf("%s waits\n", name)
		default:
			fmt.Printf("The capacity is full so %s leaves\n", name)
		}
	} else {
		fmt.Printf("The shop is closed, so %s leaves\n", name)
	}
}
