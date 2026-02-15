package user

import (
	"regexp"

	"nh-be/internal/utils/httputil"
)

const (
	usernameMinLength    = 3
	usernameMaxLength    = 30
	usernameAllowedChars = `^[a-zA-Z0-9._]+$`
)

// ValidateUsername validates a username against all rules:
// - Length: 3-30 characters
// - Allowed chars: a-z, A-Z, 0-9, _, .
// - Must start with a letter
// - Must end with a letter or number
// - No consecutive dots or underscores
// - No adjacent dot/underscore combinations
// - Not a reserved name
func ValidateUsername(username string) error {
	// Check length
	if len(username) < usernameMinLength || len(username) > usernameMaxLength {
		return ErrInvalidUsernameLength
	}

	// Check allowed characters
	if !regexp.MustCompile(usernameAllowedChars).MatchString(username) {
		return ErrInvalidUsernameChars
	}

	// Must start with a letter
	if !httputil.ValidateStartLetter(username, []string{httputil.IS_LETTER}) {
		return ErrUsernameMustStartWithLetter
	}

	// Must end with a letter or number
	if !httputil.ValidateEndChar(username, []string{httputil.IS_LETTER, httputil.IS_NUMBER}) {
		return ErrUsernameMustEndWithLetterOrNumber
	}

	// No consecutive dots or underscores
	if !httputil.ValidateNoConsecutive(username) {
		return ErrUsernameNoConsecutiveSpecialChars
	}

	// No adjacent dot/underscore combinations
	if !httputil.ValidateNoAdjacentSpecialChars(username) {
		return ErrUsernameNoAdjacentSpecialChars
	}

	// Check reserved names
	if !httputil.ValidateNotReservedName(username) {
		return ErrReservedUsername
	}

	return nil
}
