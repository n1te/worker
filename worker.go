package worker

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

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

func (w Worker) Run(job Job, config *Config) {
	c := Consumer{}
	var err error
	deliveries, err := c.Connect(config)
	if err != nil {
		log.Fatalf(" [*] Connection error: %s, exiting...", err)
	}

	log.Print(" [*] Waiting for messages. To exit press CTRL+C")

	for {
		d, ok := <-deliveries
		if ok {
			go w.run_job(d, job)
		} else {
			log.Printf(" [*] Connection closed. Trying to reconnect in %d seconds...", config.ReconnectTime)
			for {
				time.Sleep(time.Second * time.Duration(config.ReconnectTime))
				deliveries, err = c.Connect(config)
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
