package main

import "fmt"

func main() {
	s := []int{1, 2, 3} // cap 3
	fmt.Printf("ptr=%p len=%d, cap=%d\n", &s[0], len(s), cap(s))

	fmt.Println("Before modification:", s) // 1 2 3
	modifyElement(s)
	fmt.Println("After modification:", s) // 1 5 3 backed array is the same
	fmt.Printf("ptr=%p len=%d, cap=%d\n", &s[0], len(s), cap(s))

	fmt.Println("Before addition:", s) // 1 5 3
	addElement(s)
	fmt.Println("After addition:", s) // 1 5 3

	fmt.Printf("ptr=%p len=%d, cap=%d\n", &s[0], len(s), cap(s))
}

func addElement(s []int) {
	s = append(s, 10) // backed array is changed
	s[0] = 15

	fmt.Printf("ptr=%p len=%d, cap=%d\n", &s[0], len(s), cap(s))

	fmt.Println("After inner addition:", s) // 15 5 3 10
}

func modifyElement(s []int) {
	s[1] = 5
}
