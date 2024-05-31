package main

import (
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

func worker(jobs <-chan string) {
	for id := range jobs {
		log.Println("Working on image with ID", strings.Split(id, "-")[0])
		cmd := exec.Command("./thread-art-rust.exe", filepath.Join(
			args.DataPath, id, "in.png"),
			filepath.Join(args.DataPath, id, "out.png"))
		_ = cmd.Run()
	}
}

func workerPool() chan string {
	jobs := make(chan string, 100)
	for w := 1; w <= 3; w++ {
		go worker(jobs)
	}
	return jobs
}
