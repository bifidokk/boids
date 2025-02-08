package main

import "fmt"

func main() {
	array := [3]int{1, 2, 3}
	slice := array[:]

	fmt.Printf("%p \n", &array)
	fmt.Printf("%p \n", &slice[0])

	fmt.Println("Array before modifying:", array) // 123
	modifyArray(array)
	fmt.Println("Array after modifying:", array) // 123 - because we pass by copy

	fmt.Println("Slice before modifying:", slice) // 123
	modifySlice(slice)
	fmt.Println("Slice after modifying:", slice) // 100 2 3 - pass by copy but the backing array is the same
	fmt.Println("Final array:", array)           // 100 2 3 - baking array was changed by modifySlice
}

func modifyArray(array [3]int) {
	array[0] = 100
	fmt.Println("Array inside modifying:", array) // 100, 2, 3
}

func modifySlice(slice []int) {
	fmt.Printf("%p \n", &slice)    // different address because passed by copy
	fmt.Printf("%p \n", &slice[0]) // same address because internal pointer copied but it still points to the same address

	slice[0] = 100
	fmt.Println("Slice inside modifying:", slice) // 100, 2, 3
}
