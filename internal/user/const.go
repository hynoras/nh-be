package user

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateUsername = errors.New("username already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrForbidCreateUser  = errors.New("you do not have permission to create users")
	ErrForbidViewUsers   = errors.New("you do not have permission to view users")
	ErrForbidViewUser    = errors.New("you do not have permission to view this user")
	ErrForbidUpdateUser  = errors.New("you do not have permission to update this user")
	ErrForbidDeleteUser  = errors.New("you do not have permission to delete this user")
)
