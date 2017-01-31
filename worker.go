package worker

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/streadway/amqp"
)

var conf Config

type Job func([]byte) bool

type Worker struct {
}

func (w Worker) run_job(d amqp.Delivery, job Job) {
	log.Printf("Received a message: %s", d.Body)
	result := job(d.Body)
	log.Print("Done")
	if result {
		d.Ack(false)
	} else {
		d.Reject(true)
	}
}

func (w Worker) Run(job Job) {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	conf = LoadConfig(dir + "/config.toml")

	c := Consumer{}
	var err error
	deliveries, err := c.Connect(conf)
	if err != nil {
		log.Fatalf(" [*] Connection error: %s, exiting...", err)
	}

	log.Print(" [*] Waiting for messages. To exit press CTRL+C")

	for {
		d, ok := <-deliveries
		if ok {
			go w.run_job(d, job)
		} else {
			log.Printf(" [*] Connection closed. Trying to reconnect in %d seconds...", conf.ReconnectTime)
			for {
				time.Sleep(time.Second * time.Duration(conf.ReconnectTime))
				deliveries, err = c.Connect(conf)
				if err != nil {
					log.Printf(" [*] Reconnect failed: %s", err)
				} else {
					log.Print(" [*] Reconnected")
					log.Print(" [*] Waiting for messages. To exit press CTRL+C")
					break
				}
			}
		}
	}
}
