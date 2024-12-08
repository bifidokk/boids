package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	matches   []string
	waitgroup = sync.WaitGroup{}
	lock      = sync.Mutex{}
)

func main() {
	waitgroup.Add(1)
	go fileSearch("/Users/danil/repo", "README.md")
	waitgroup.Wait()

	for _, match := range matches {
		fmt.Println("Matched: ", match)
	}
}

func fileSearch(root string, filename string) {
	fmt.Println("Searching: ", root, filename)

	files, _ := os.ReadDir(root)

	for _, file := range files {
		if file.Name() == filename {
			lock.Lock()
			matches = append(matches, filepath.Join(root, file.Name()))
			lock.Unlock()
		}

		if file.IsDir() {
			waitgroup.Add(1)
			go fileSearch(filepath.Join(root, file.Name()), filename)
		}
	}

	waitgroup.Done()
}
