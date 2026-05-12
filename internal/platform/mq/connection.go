package mq

import (
	"fmt"
	"log/slog"
	"nh-be/internal/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRabbitMQConnection(cfg *config.Config) (*amqp.Connection, error) {
	conn, err := amqp.Dial(
		fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.RabbitMQUsername, cfg.RabbitMQPassword, cfg.RabbitMQHost, cfg.RabbitMQPort),
	)
	if err != nil {
		return nil, err
	}

	slog.Info("Connected to RabbitMQ")
	return conn, nil
}
