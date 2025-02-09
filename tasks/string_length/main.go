package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	s := "abЯЮ漢"

	// len returns count of bytes
	fmt.Println(len(s))                    // 9 a, b - 2 bytes, ЯЮ - 4 bytes, 漢 - 3 bytes
	fmt.Println(utf8.RuneCountInString(s)) // 5
}
