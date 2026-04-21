package mq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Publish(ctx context.Context, ch *amqp.Channel, exchange, key, body string, mandatory, immediate bool) error {
	return ch.PublishWithContext(
		ctx,
		exchange,
		key,
		mandatory,
		immediate,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		},
	)
}
