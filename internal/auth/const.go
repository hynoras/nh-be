package auth

import "errors"

var (
	ErrVerficationTokenNotFound = errors.New("token not found")
	ErrUnauthenticated          = errors.New("unauthenticated")
	ErrEmailAlreadyExists       = errors.New("email already exists")
)

const (
	AuthExchangeName           = "auth"
	UserRegisteredRoutingKey   = "user.registered"
	SendVerificationEmailQueue = "send-verification-email"
)
