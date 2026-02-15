package email

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"nh-be/mq"
	"nh-be/utils"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EmailConsumer interface {
	SendVerificationEmail(ctx context.Context) error
}

type emailConsumer struct {
	channel *amqp.Channel
}

func NewEmailConsumer(ch *amqp.Channel) EmailConsumer {
	return &emailConsumer{channel: ch}
}

func (s *emailConsumer) SendVerificationEmail(ctx context.Context) error {
	resendClient := NewResendClient()
	frontendURL := utils.MustEnv("FRONTEND_URL")
	verifyEmailSuffixURL := utils.MustEnv("VERIFY_EMAIL_SUFFIX_URL")
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
				d.Nack(false, false)
				continue
			}

			verificationURL := fmt.Sprintf("%s%s?token=%s", frontendURL, verifyEmailSuffixURL, req.Token)

			htmlContent, htmlErr := utils.ConvertHtmlToString("verification_email.html", map[string]string{
				"UserName":        utils.ExtractUsernameFromEmail(req.ToEmail),
				"VerificationURL": verificationURL,
			})

			if htmlErr != nil {
				log.Printf("Failed to convert HTML to string: %v", htmlErr)
				d.Nack(false, false)
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
				d.Nack(false, false)
				continue
			}

			d.Ack(false)
		}
	}()
	<-ctx.Done()
	return ctx.Err()
}
