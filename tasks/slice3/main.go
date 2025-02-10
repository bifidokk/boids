package main

import (
	"fmt"
	"math/rand"
)

func main() {
	fmt.Println(uniqN(99))
}

func uniqN(n int) []int {
	m := make(map[int]struct{}, n)
	s := make([]int, 0, n)

	for len(s) < n {
		v := rand.Intn(100)

		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			s = append(s, v)
		}
	}

	return s
}
