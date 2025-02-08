package main

import "fmt"

type car struct {
	color   string
	mileage int
}

func main() {
	// len 3, capacity 3
	cars := []car{
		{color: "red", mileage: 5},
		{color: "blue", mileage: 4},
		{color: "green", mileage: 3},
	}

	fmt.Printf("%d %d %p \n", len(cars), cap(cars), cars)

	carPointer := &cars[0]
	carPointer.mileage += 2

	// because capacity is 3, we change backing array
	cars = append(cars, car{color: "yellow", mileage: 8})
	// after append len 4, capacity 6 and address is changed. So now pointer doesn't point on cars[0]
	fmt.Printf("%d %d %p \n", len(cars), cap(cars), cars)

	carPointer.mileage += 2

	fmt.Println(cars[0].mileage, cars[0].color)       // 7, red
	fmt.Println(carPointer.mileage, carPointer.color) // 9, red
}
