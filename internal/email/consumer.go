package email

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"nh-be/internal/config"
	"nh-be/internal/utils/stringutil"
	"nh-be/mq"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EmailConsumer interface {
	SendVerificationEmail(ctx context.Context) error
}

type emailConsumer struct {
	channel              *amqp.Channel
	frontendURL          string
	verifyEmailSuffixURL string
	resendAPIKey         string
}

func NewEmailConsumer(ch *amqp.Channel, cfg *config.Config) EmailConsumer {
	return &emailConsumer{
		channel:              ch,
		frontendURL:          cfg.FrontendURL,
		verifyEmailSuffixURL: cfg.VerifyEmailSuffixURL,
		resendAPIKey:         cfg.ResendAPIKey,
	}
}

func (s *emailConsumer) SendVerificationEmail(ctx context.Context) error {
	resendClient := NewResendClient(s.resendAPIKey)
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
			var req SendVerificationEmailDto

			if err := json.Unmarshal(d.Body, &req); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				_ = d.Nack(false, false)
				continue
			}

			verificationURL := fmt.Sprintf("%s%s?token=%s", s.frontendURL, s.verifyEmailSuffixURL, req.Token)

			htmlContent, htmlErr := ConvertHtmlToString("verification_email.html", map[string]string{
				"UserName":        stringutil.ExtractUsernameFromEmail(req.ToEmail),
				"VerificationURL": verificationURL,
			})

			if htmlErr != nil {
				log.Printf("Failed to convert HTML to string: %v", htmlErr)
				_ = d.Nack(false, false)
				continue
			}

			sendEmailErr := SendEmail(
				resendClient,
				"Acme <onboarding@resend.dev>",
				req.ToEmail,
				"Verify your email",
				htmlContent,
			)

			if sendEmailErr != nil {
				log.Printf("Failed to send email: %v", sendEmailErr)
				_ = d.Nack(false, false)
				continue
			}

			_ = d.Ack(false)
		}
	}()
	<-ctx.Done()
	return ctx.Err()
}
