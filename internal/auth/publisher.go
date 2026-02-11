package auth

import (
	"context"
	"encoding/json"
	"log"
	"nh-be/mq"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AuthPublisher interface {
	PublishSendVerificationEmail(ctx context.Context, email string) error
}

type authPublisher struct {
	channel *amqp.Channel
}

func NewAuthPublisher(ch *amqp.Channel) AuthPublisher {
	return &authPublisher{channel: ch}
}

func (s *authPublisher) PublishSendVerificationEmail(ctx context.Context, email string) error {
	body, err := json.Marshal(map[string]string{
		"email": email,
	})
	if err != nil {
		return err
	}

	err = mq.Publish(ctx, s.channel, AuthExchangeName, SendVerificationEmailQueue, string(body), false, false)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
		return err
	}
	return nil
}
