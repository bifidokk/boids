package main

import "fmt"

func main() {
	s := "G🤔o" // 🤔- 4 bytes

	for i := 0; i < len(s); i++ {
		fmt.Printf("%c\n", s[i])
	} // G****o
}
