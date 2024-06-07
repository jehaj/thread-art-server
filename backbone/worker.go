package main

import (
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

// worker is a worker thread that runs the binary on the image with the given ID. The file structure is
// [datapath]/[id]/[in.png] and outputs [datapath]/[id]/[out.png]. It marks the image as finished when the binary
// is done running.
func worker(s *Service, jobs <-chan string) {
	for id := range jobs {
		log.Println("Working on image with ID", strings.Split(id, "-")[0])
		cmd := exec.Command("./thread-art-rust.exe", filepath.Join(
			args.DataPath, id, "in.png"),
			filepath.Join(args.DataPath, id, "out.png"))
		_ = cmd.Run()
		s.ImageFinished(id)
	}
}

// workerPool returns a channel. It has three worker threads. The service is needed such that the worker
// thread can mark the image as finished when done.
func workerPool(s *Service) chan string {
	jobs := make(chan string, 100)
	for w := 1; w <= 3; w++ {
		go worker(s, jobs)
	}
	return jobs
}
