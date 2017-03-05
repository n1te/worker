package main

import (
	"log"
	"path/filepath"
	"os"

	"github.com/n1te/worker"
)

type MyConfig struct {
	worker.Config
	TestString string `toml:"test_string"`
}

var c MyConfig

func job(data []byte) bool {
	log.Printf("Received a message: %s", data)
	return true
}

func main() {
	var w worker.Worker
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	worker.LoadConfigFromFile(&c, dir + "/config.toml")
	w.Run(job, &c.Config)
}
