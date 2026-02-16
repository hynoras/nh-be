package auth

import "errors"

var (
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrVerificationTokenNotFound = errors.New("token not found")
	ErrUnauthenticated           = errors.New("unauthenticated")
	ErrEmailAlreadyExists        = errors.New("email already exists")
	ErrVerificationTokenExpired  = errors.New("token expired")
	ErrInvalidVerificationToken  = errors.New("invalid token, please check token type")
	ErrSessionNotFound           = errors.New("session not found")
)
