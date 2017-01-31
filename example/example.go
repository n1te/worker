package main

import (
	"github.com/n1te/worker"
	"log"
)

func job(data []byte) bool {
	log.Print("Hello WRLD")
	return true
}

func main() {
	var w worker.Worker
	w.Run(job)
}
