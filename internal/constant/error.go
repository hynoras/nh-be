package constant

import "errors"

const (
	ErrAuthorizationFailed = "Authorization failed"
)

var (
	ErrInvalidIDFormat = errors.New("invalid id format")
)
