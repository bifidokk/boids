package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"reflect"
	"strings"
)

func main() {
	f1 := `aaa
ddd
`
	f2 := `bbb
eee
`
	f3 := `ccc
fff
`
	files := []io.Reader{
		strings.NewReader(f1),
		strings.NewReader(f2),
		strings.NewReader(f3),
	}

	rows, err := ConcurrentSortHead(4, files...)
	if err != nil {
		panic(err)
	}

	fmt.Println(rows)
	if !reflect.DeepEqual(rows, []string{"aaa", "bbb", "ccc", "ddd"}) {
		panic("wrong code")
	}

}

func ConcurrentSortHead(m int, files ...io.Reader) ([]string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chans := make(map[chan string]string)

	for _, f := range files {
		ch := make(chan string)
		chans[ch] = ""

		go func(f io.Reader, c chan string) {
			defer close(c)

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				select {
				case c <- scanner.Text():
				case <-ctx.Done():
				}
			}

			if err := scanner.Err(); err != nil {
				panic(err)
			}
		}(f, ch)
	}

	for ch := range chans {
		row, ok := <-ch
		if !ok {
			delete(chans, ch)
			continue
		}
		chans[ch] = row
	}

	result := make([]string, 0, m)

	for len(chans) > 0 && len(result) < m {
		var minCh chan string
		var minRow string

		firstIter := true
		for ch, row := range chans {
			if firstIter {
				minRow, minCh = row, ch
				firstIter = false
				continue
			}

			if row < minRow {
				minRow, minCh = row, ch
			}
		}

		result = append(result, minRow)
		row, ok := <-minCh

		if !ok {
			delete(chans, minCh)
			continue
		}

		chans[minCh] = row
	}

	return result, nil
}
