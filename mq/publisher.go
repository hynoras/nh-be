package mq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Publish(ctx context.Context, ch *amqp.Channel, exchange, name string, body string, mandatory, immediate bool) error {
	return ch.PublishWithContext(
		ctx,
		exchange,
		name,
		mandatory,
		immediate,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		},
	)
}
