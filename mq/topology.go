package mq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func DeclareQueues(ch *amqp.Channel, name string) error {
	_, err := ch.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
	return err
}

func DeclareExchange(ch *amqp.Channel, name string) error {
	err := ch.ExchangeDeclare(
		name,
		"direct",
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
	return err
}

func BindQueue(ch *amqp.Channel, queueName, routingKey, exchangeName string) error {
	return ch.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
}
