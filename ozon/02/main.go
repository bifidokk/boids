package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
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
		nString, _ := reader.ReadString('\n')
		n, _ := strconv.Atoi(strings.TrimSpace(nString))

		inputLine, _ := reader.ReadString('\n')
		inputLine = strings.ReplaceAll(inputLine, "\n", "")
		input := strings.Fields(inputLine)

		outputLine, _ := reader.ReadString('\n')
		outputLine = strings.ReplaceAll(outputLine, "\n", "")
		output := strings.Fields(outputLine)

		if strings.TrimSpace(inputLine) != inputLine || strings.TrimSpace(outputLine) != outputLine {
			fmt.Fprintln(writer, "no")
			continue
		}

		if len(input) != n || len(output) != n {
			fmt.Fprintln(writer, "no")
			continue
		}

		if !isValidFormat(inputLine) || !isValidFormat(outputLine) {
			fmt.Fprintln(writer, "no")
			continue
		}

		inputNumbers := make([]int, n)
		validInput := true
		for j, numStr := range input {
			if !isValidNumber(numStr) {
				validInput = false
				break
			}

			num, err := strconv.Atoi(numStr)
			if err != nil || num < -1e9 || num > 1e9 {
				validInput = false
				break
			}
			inputNumbers[j] = num
		}

		if !validInput {
			fmt.Fprintln(writer, "no")
			continue
		}

		outputNumbers := make([]int, n)
		validOutput := true
		for j, numStr := range output {
			if !isValidNumber(numStr) {
				validInput = false
				break
			}

			num, err := strconv.Atoi(numStr)
			if err != nil || num < -1e9 || num > 1e9 {
				validOutput = false
				break
			}
			outputNumbers[j] = num
		}

		if !validOutput {
			fmt.Fprintln(writer, "no")
			continue
		}

		sort.Ints(inputNumbers)
		isSorted := true
		for j := 0; j < n; j++ {
			if inputNumbers[j] != outputNumbers[j] {
				isSorted = false
				break
			}
		}

		if isSorted {
			fmt.Fprintln(writer, "yes")
		} else {
			fmt.Fprintln(writer, "no")
		}
	}
}

func isValidNumber(s string) bool {
	if len(s) == 0 {
		return false
	}

	if s[0] == '-' {
		s = s[1:]
	}

	if len(s) > 1 && s[0] == '0' {
		return false
	}

	_, err := strconv.Atoi(s)
	return err == nil
}

func isValidFormat(line string) bool {
	re := regexp.MustCompile(`^-?\d+(\s-?\d+)*$`)
	return re.MatchString(line)
}
