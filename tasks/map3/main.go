package main

import "fmt"

type WordCounter struct {
	m     map[string]int
	order []string
	limit int
}

func NewWordCounter(limit int) *WordCounter {
	return &WordCounter{
		m:     make(map[string]int),
		limit: limit,
	}
}

func (w *WordCounter) CountWord(word string) {
	if _, ok := w.m[word]; !ok {
		w.order = append(w.order, word)
	}

	w.m[word]++

	if len(w.m) > w.limit {
		delete(w.m, w.order[0])
		w.order = w.order[1:]
	}
}

func main() {
	w := NewWordCounter(3)

	words := []string{"orange", "apple", "pear", "grapefruit", "banana", "lemon", "pear"}

	for _, word := range words {
		w.CountWord(word)
	}

	fmt.Println(w.m)
}
