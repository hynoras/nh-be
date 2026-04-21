package email

type SendVerificationEmailDto struct {
	ToEmail string
	Token   string
}
