package utils

import (
	"github.com/resend/resend-go/v3"
)

func NewResendClient() *resend.Client {
	apiKey := MustEnv("RESEND_API_KEY")
	return resend.NewClient(apiKey)
}

func SendEmail(client *resend.Client, from, to, subject, htmlContent string) error {

	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{to},
		Subject: subject,
		Html:    htmlContent,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		return err
	}

	return nil
}
