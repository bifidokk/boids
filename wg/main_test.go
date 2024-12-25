package main

import (
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func Test_updateMessage(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go updateMessage("updated", &wg)
	wg.Wait()

	if msg != "updated" {
		t.Errorf("updateMessage failed: expected \"updated\", got \"%s\"", msg)
	}
}

func Test_printMessage(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	msg = "updated"
	printMessage()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, msg) {
		t.Errorf("printMessage failed: expected to contain \"updated\", got \"%s\"", output)
	}
}

func Test_main(t *testing.T) {
	stdOut := os.Stdout

	r, w, _ := os.Pipe()
	os.Stdout = w

	main()

	_ = w.Close()
	result, _ := io.ReadAll(r)
	output := string(result)
	os.Stdout = stdOut

	if !strings.Contains(output, "Hello, universe!") {
		t.Errorf("main failed: expected to contain \"Hello, universe!\", got \"%s\"", output)
	}

	if !strings.Contains(output, "Hello, cosmos!") {
		t.Errorf("main failed: expected to contain \"Hello, cosmos!\", got \"%s\"", output)
	}

	if !strings.Contains(output, "Hello, world!") {
		t.Errorf("main failed: expected to contain \"Hello, world!\", got \"%s\"", output)
	}
}
