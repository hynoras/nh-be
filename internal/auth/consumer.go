package auth

import (
	"context"
	"encoding/json"
	"log"
	"nh-be/mq"
	"nh-be/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AuthConsumer interface {
	ConsumeSendVerificationEmail(ctx context.Context) error
}

type authConsumer struct {
	channel *amqp.Channel
}

func NewAuthConsumer(ch *amqp.Channel) AuthConsumer {
	return &authConsumer{channel: ch}
}

func (s *authConsumer) ConsumeSendVerificationEmail(ctx context.Context) error {
	resendClient := utils.NewResendClient()
	msgs, err := mq.Consumer(
		ctx,
		s.channel,
		amqp.Queue{
			Name: SendVerificationEmailQueue,
		},
	)

	if err != nil {
		log.Printf("Failed to start consumer: %v", err)
		return err
	}

	go func() {
		for d := range msgs {
			var req map[string]string
			if err := json.Unmarshal(d.Body, &req); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				d.Nack(false, false)
				continue
			}
			err = utils.SendEmail(
				resendClient,
				"Acme <onboarding@resend.dev>",
				req["email"],
				"Verify your email",
				"Please verify your email",
			)
			if err != nil {
				log.Printf("Failed to send email: %v", err)
				d.Nack(false, false)
				continue
			}
			d.Ack(false)
		}
	}()
	<-ctx.Done()
	return ctx.Err()
}
