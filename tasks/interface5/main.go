package main

import "fmt"

type SomeStruct struct {
	Value int
}

func checkForNil(i interface{}) {
	if i == nil {
		fmt.Println("it's nil!")
		return
	}

	fmt.Println("it's not nil!")
}

func main() {
	var s *SomeStruct // prints not nil because s holds type information (SomeStruct)
	var k interface{}
	checkForNil(s)
	checkForNil(k)
}
