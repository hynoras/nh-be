package mq

import (
	"fmt"
	"log"
	"nh-be/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQConnection() (*amqp.Connection, error) {
	host := utils.MustEnv("RABBITMQ_HOST")
	port := utils.MustEnvInt("RABBITMQ_PORT")
	user := utils.MustEnv("RABBITMQ_USERNAME")
	pass := utils.MustEnv("RABBITMQ_PASSWORD")

	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%d/", user, pass, host, port),
	)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to RabbitMQ")
	return conn, nil
}
