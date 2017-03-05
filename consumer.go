package worker

import "github.com/streadway/amqp"

type Consumer struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue
}

func (c Consumer) Connect(conf *Config) (<-chan amqp.Delivery, error) {
	var err error

	c.conn, err = amqp.Dial(conf.Rabbit.GetConnectionString())
	if err != nil {
		return nil, err
	}

	c.ch, err = c.conn.Channel()
	if err != nil {
		return nil, err
	}

	c.q, err = c.ch.QueueDeclare(
		conf.QueueName, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		return nil, err
	}

	err = c.ch.Qos(conf.PrefetchCount, 0, false)

	var deliveries <-chan amqp.Delivery
	deliveries, err = c.ch.Consume(
		c.q.Name, // queue
		"",       // consumer
		false,    // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		return nil, err
	}

	return deliveries, nil
}
