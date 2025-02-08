package main

import "fmt"

func main() {
	x := 1
	y := 2

	defer func(val int) {
		fmt.Println("x:", val)
	}(x) // 1, because we pass current value

	defer func() {
		fmt.Println("y:", y)
	}() // 200 because y value is changed

	x = 100
	y = 200
}
