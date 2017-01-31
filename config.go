package worker

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

type Rabbit struct {
	User     string
	Password string
	Host     string
	Port     int
}

func (r Rabbit) GetConnectionString() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", r.User, r.Password, r.Host, r.Port)
}

type Config struct {
	QueueName     string `toml:"queue_name"`
	ReconnectTime int    `toml:"reconnect_time"`
	PrefetchCount int    `toml:"prefetch_count"`
	Rabbit        Rabbit
}

func LoadConfig(configFile string) Config {

	config := Config{}
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}

	return config
}
