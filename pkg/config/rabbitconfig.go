package config

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type Configuration struct {
	Host             string
	Port             string
	Username         string
	Password         string
	ConnectionPrefix string
}

func NewConfiguration() *Configuration {
	host := getOrDefault("RABBITMQ_HOST", "localhost")
	port := getOrDefault("RABBITMQ_PORT", "5672")
	username := getOrDefault("RABBITMQ_USERNAME", "test")
	password := getOrDefault("RABBITMQ_PASSWORD", "test")
	applicationName := getOrDefault("APPLICATION_NAME", "go-application")

	return &Configuration{
		Host:             host,
		Port:             port,
		Username:         username,
		Password:         password,
		ConnectionPrefix: applicationName,
	}
}

func (config *Configuration) Connect() *Connection {

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Username, config.Password, config.Host, config.Port)

	amqpConfig := amqp091.Config{Properties: map[string]interface{}{"connection_name": config.ConnectionPrefix}}

	conn, err := amqp091.DialConfig(url, amqpConfig)

	if err != nil {
		log.Panicf("%s: %s", "Failed to connect to RabbitMQ", err)
	}

	return &Connection{
		Connection: conn,
	}
}

func openChannel(conn *amqp091.Connection) *amqp091.Channel {
	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("%s: %s", "Failed to open a channel", err)
	}
	return ch
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

//func getOrDefault(envVar string, defaultVal string) string {
//	val := os.Getenv(envVar)
//	if val == "" {
//		return defaultVal
//	}
//	return val
//}
