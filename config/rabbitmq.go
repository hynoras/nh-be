package config

import (
	"fmt"
	"log"
	"nh-be/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectRabbitMQ() {
	host := utils.MustEnv("RABBITMQ_HOST")
	port := utils.MustEnvInt("RABBITMQ_PORT")
	user := utils.MustEnv("RABBITMQ_USERNAME")
	pass := utils.MustEnv("RABBITMQ_PASSWORD")

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", user, pass, host, port))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	log.Println("Successfully connected to RabbitMQ")
}
