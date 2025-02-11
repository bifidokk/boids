package main

import "fmt"

func main() {
	data := make([]int, 5, 10)

	for i := range data {
		data[i] = i + 1
	}

	fmt.Printf("Initial %v ptr=%p len=%d, cap=%d\n", data, &data[0], len(data), cap(data)) // 1 2 3 4 5, len 5, cap 10

	data = sliceCapacity(data, 1)

	fmt.Printf("After modification 1 %v ptr=%p len=%d, cap=%d\n", data, &data[0], len(data), cap(data)) // 2 3 4 5, len 4, cap 9

	newData := make([]int, 0, 3)
	newData = sliceCapacity(data, 2)

	fmt.Printf("After modification 2 %v ptr=%p len=%d, cap=%d\n", newData, &newData[0], len(newData), cap(newData)) // 4, 5, cap 7
}

func sliceCapacity(data []int, start int) []int {
	subSlice := data[start:]

	return subSlice
}
