package email

func MapToSendVerificationEmailDto(toEmail string, rawToken string) SendVerificationEmailDto {
	return SendVerificationEmailDto{
		ToEmail: toEmail,
		Token:   rawToken,
	}
}
