package user

import (
	"errors"
	"net/http"
	"nh-be/internal/utils/httputil"
)

const (
	ErrGetAllUsersFailed   = "Failed to get users"
	ErrGetUserDetailFailed = "Failed to get user detail"
	ErrCreateUserFailed    = "Failed to create user"
	ErrUpdateUserFailed    = "Failed to update user"
	ErrDeleteUserFailed    = "Failed to delete user"
)

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

	// Proactive named errors
	ErrAssignedUserNotFound = errors.New("assigned users not found")
	ErrUsersNotFound        = errors.New("users not found")
)

func init() {
	httputil.RegisterError(ErrUserNotFound, http.StatusNotFound, "User not found")
	httputil.RegisterError(ErrForbidViewUsers, http.StatusForbidden, ErrGetAllUsersFailed)
	httputil.RegisterError(ErrForbidViewUser, http.StatusForbidden, ErrGetUserDetailFailed)
	httputil.RegisterError(ErrForbidUpdateUser, http.StatusForbidden, ErrUpdateUserFailed)
	httputil.RegisterError(ErrForbidDeleteUser, http.StatusForbidden, ErrDeleteUserFailed)
	httputil.RegisterError(ErrDuplicateUsername, http.StatusConflict, "Invalid username")
	httputil.RegisterError(ErrDuplicateEmail, http.StatusConflict, "Invalid email")
	httputil.RegisterError(ErrUsernameMustStartWithLetter, http.StatusBadRequest, "Invalid username")
	httputil.RegisterError(ErrUsernameMustEndWithLetterOrNumber, http.StatusBadRequest, "Invalid username")
	httputil.RegisterError(ErrUsernameNoConsecutiveSpecialChars, http.StatusBadRequest, "Invalid username")
	httputil.RegisterError(ErrUsernameNoAdjacentSpecialChars, http.StatusBadRequest, "Invalid username")
	httputil.RegisterError(ErrReservedUsername, http.StatusBadRequest, "Invalid username")
}
