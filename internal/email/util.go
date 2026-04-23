package email

import (
	"bytes"
	"embed"
	"text/template"

	"github.com/resend/resend-go/v3"
)

var templateFS embed.FS
var tmpl = template.Must(template.ParseFS(templateFS, "templates/*.html"))

func NewResendClient(apiKey string) *resend.Client {
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

func ConvertHtmlToString(name string, data any) (string, error) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, name, data)
	return buf.String(), err
}
