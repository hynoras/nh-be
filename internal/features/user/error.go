package user

import "errors"

var (
	ErrUserNotFound                      = errors.New("user not found")
	ErrDuplicateUsername                 = errors.New("username already exists")
	ErrDuplicateEmail                    = errors.New("email already exists")
	ErrForbidCreateUser                  = errors.New("you do not have permission to create users")
	ErrForbidViewUsers                   = errors.New("you do not have permission to view users")
	ErrForbidViewUser                    = errors.New("you do not have permission to view this user")
	ErrForbidUpdateUser                  = errors.New("you do not have permission to update this user")
	ErrForbidDeleteUser                  = errors.New("you do not have permission to delete this user")
	ErrInvalidUsernameLength             = errors.New("username must be between 3 and 30 characters")
	ErrInvalidUsernameChars              = errors.New("username can only contain letters, numbers, dots, and underscores")
	ErrUsernameMustStartWithLetter       = errors.New("username must start with a letter")
	ErrUsernameMustEndWithLetterOrNumber = errors.New("username must end with a letter or number")
	ErrUsernameNoConsecutiveSpecialChars = errors.New("username cannot contain consecutive dots or underscores")
	ErrUsernameNoAdjacentSpecialChars    = errors.New("username cannot contain adjacent dots and underscores")
	ErrReservedUsername                  = errors.New("username is reserved and cannot be used")
)
