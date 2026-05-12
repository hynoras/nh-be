package email

import (
	"context"
	"encoding/json"
	"log/slog"
	"nh-be/internal/platform/mq"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EmailPublisher interface {
	SendVerificationEmail(ctx context.Context, dto SendVerificationEmailDto) error
}

type emailPublisher struct {
	channel *amqp.Channel
}

func NewEmailPublisher(ch *amqp.Channel) EmailPublisher {
	return &emailPublisher{channel: ch}
}

func (s *emailPublisher) SendVerificationEmail(ctx context.Context, dto SendVerificationEmailDto) error {
	body, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	err = mq.Publish(ctx, s.channel, AuthExchangeName, UserRegisteredRoutingKey, string(body), false, false)
	if err != nil {
		slog.Error("Failed to publish message", "error", err)
		return err
	}
	return nil
}
