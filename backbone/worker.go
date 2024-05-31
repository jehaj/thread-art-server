package main

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan string) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
	}
}

func workerPool() chan string {
	jobs := make(chan string, 100)
	for w := 1; w <= 3; w++ {
		go worker(w, jobs)
	}
	return jobs
}
