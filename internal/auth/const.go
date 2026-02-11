package auth

import "errors"

var ErrUnauthenticated = errors.New("unauthenticated")

const (
	AuthExchangeName           = "auth"
	UserRegisteredRoutingKey   = "user.registered"
	SendVerificationEmailQueue = "send-verification-email"
)
