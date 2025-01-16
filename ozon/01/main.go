package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	line, _ := reader.ReadString('\n')
	t, _ := strconv.Atoi(strings.TrimSpace(line))

	for i := 0; i < t; i++ {
		line, _ = reader.ReadString('\n')
		salaryString := strings.TrimSpace(line)

		if len(salaryString) <= 1 {
			fmt.Fprintln(writer, 0)
			continue
		}

		// If remove the first "out of order" digit (the one that is smaller than the next), the number will remain the maximum.
		indexToRemove := -1
		for j := 0; j < len(salaryString)-1; j++ {
			if salaryString[j] < salaryString[j+1] {
				indexToRemove = j
				break
			}
		}

		if indexToRemove == -1 {
			indexToRemove = len(salaryString) - 1
		}

		result := salaryString[:indexToRemove] + salaryString[indexToRemove+1:]
		fmt.Fprintln(writer, result)
	}
}

// Input
//3 - number of test cases
//9
//0
//9123

// Output
//0
//0
//923 - max value if one of digits removed
