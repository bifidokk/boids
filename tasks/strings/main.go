package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {
	str := ""

	start := time.Now()
	for i := 0; i < 100_000; i++ {
		str += fmt.Sprintf("%d", i) // string is immutable. During concatenation, we create a new instance of string and allocate memory
	}

	fmt.Println("Elapsed time:", time.Since(start))

	// better approach
	start2 := time.Now()
	builder := strings.Builder{}

	for i := 0; i < 100_000; i++ {
		builder.WriteString(strconv.Itoa(i))
	}

	fmt.Println("Elapsed time:", time.Since(start2))
}
