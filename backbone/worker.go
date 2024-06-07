package main

import (
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

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

func workerPool(s *Service) chan string {
	jobs := make(chan string, 100)
	for w := 1; w <= 3; w++ {
		go worker(s, jobs)
	}
	return jobs
}
