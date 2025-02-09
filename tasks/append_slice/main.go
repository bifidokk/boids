package main

import "fmt"

func main() {
	data := []int{1, 2, 3, 4, 5}
	fmt.Println("Initial value ", data) // 1, 2, 3, 4, 5

	modify(data[:2])
	fmt.Println("Modified value ", data) // 1, 2, 6, 7, 5
}

func modify(slice []int) {
	fmt.Println(cap(slice)) // capacity is still 5, so we append only 2 numbers, then backed array wll not be changed
	// in case we add 4 numbers, backed array will be changed and the last print in main function will print 1 2 3 4 5
	slice = append(slice, 6, 7)
	fmt.Println("After modifying slice ", slice) // 1, 2, 6, 7
}
