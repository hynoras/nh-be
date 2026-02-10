package utils

import (
	"fmt"
	"log"

	"github.com/resend/resend-go/v3"
)

func SendEmail(from, to, subject, htmlContent string) {

	apiKey := MustEnv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{to},
		Subject: subject,
		Html:    htmlContent,
	}

	_, err := client.Emails.Send(params)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	fmt.Println("Email sent successfully")
}
