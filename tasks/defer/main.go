package main

import "fmt"

func main() {
	x := 1
	y := 2

	defer func(val int) {
		fmt.Println("x:", val)
	}(x) // 1, because we pass current value by copy, will be shown 2nd

	defer func() {
		fmt.Println("y:", y)
	}() // 200, because y value is changed, we pass by pointer, will be shown 1st

	x = 100
	y = 200

}
