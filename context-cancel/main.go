package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/longtask", longTaskHandler)

	server := &http.Server{
		Addr: ":8081",
	}

	fmt.Println("Server is running on http://localhost:8081")
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

func longTaskHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log.Println("Handler started")
	defer log.Println("Handler finished")

	// Simulate a long-running task
	ch := make(chan string)

	go func() {
		// Simulate processing time
		time.Sleep(5 * time.Second)
		ch <- "Task completed"
	}()

	select {
	case <-ctx.Done():
		// The request was canceled
		log.Println("Request canceled by client")
		http.Error(w, "Request canceled", http.StatusRequestTimeout)
	case result := <-ch:
		// Task completed successfully
		log.Println("Task completed")
		fmt.Fprintln(w, result)
	}
}
