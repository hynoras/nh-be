package mq

import (
	"fmt"
	"log"
	"nh-be/pkg/env"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQConnection() (*amqp.Connection, error) {
	host := env.MustEnv("RABBITMQ_HOST")
	port := env.MustEnvInt("RABBITMQ_PORT")
	user := env.MustEnv("RABBITMQ_USERNAME")
	pass := env.MustEnv("RABBITMQ_PASSWORD")

	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%d/", user, pass, host, port),
	)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to RabbitMQ")
	return conn, nil
}
