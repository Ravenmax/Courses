package main

import (
	"fmt"
	"sync"
	"time"
)

func sqrWorkers(wg *sync.WaitGroup, tasks <-chan int, results chan<- int, instance int) {
	for num := range tasks {
		time.Sleep(time.Millisecond)
		fmt.Println("[worker %v] Sending result by worker %v", instance, instance)
		results <- num * num
	}
	wg.Done()
}

func main() {
	fmt.Println("[main] Main started")

	var wg sync.WaitGroup
	tasks := make(chan int, 10)
	results := make(chan int, 10)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go sqrWorkers(&wg, tasks, results, i)
	}

	for i := 0; i < 5; i++ {
		tasks <- i * 2
	}
	fmt.Println("[main] Wrote 5 task")
	close(tasks)
	wg.Wait()
	for i := 0; i < 5; i++ {
		result := <-results
		fmt.Println("[main] Result", i, ":", result)
	}
	fmt.Println("Main stopped")
}
