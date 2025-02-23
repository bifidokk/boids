package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID       int64
	Filename string
}

func ProcessTask(task Task) string {
	time.Sleep(1 * time.Second)
	return fmt.Sprintf("Done: Task ID: %d\nFilename: %s", task.ID, task.Filename)
}

func RunWorker(id int64, taskCh <-chan Task, resultCh chan<- string) {
	for task := range taskCh {
		fmt.Printf("Worker %d received task %d\n", id, task.ID)
		resultCh <- ProcessTask(task)
		fmt.Printf("Worker %d finished task %d\n", id, task.ID)
	}
}

const numWorkers = 3
const numTasks = 10

func main() {
	taskCh := make(chan Task, numTasks)
	resultCh := make(chan string, numTasks)

	wg := sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			RunWorker(int64(i), taskCh, resultCh)
		}()
	}

	go func() {
		for i := 0; i < numTasks; i++ {
			taskCh <- Task{
				ID:       int64(i),
				Filename: fmt.Sprintf("task-%d", i),
			}
		}

		close(taskCh)
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for res := range resultCh {
		fmt.Println(res)
	}

	fmt.Println("DONE")

}
