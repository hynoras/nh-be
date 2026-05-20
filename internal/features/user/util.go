package user

import (
	"regexp"

	"nh-be/internal/utils/validationutil"
)

const (
	usernameMinLength    = 3
	usernameMaxLength    = 30
	usernameAllowedChars = `^[a-zA-Z0-9._]+$`
)

var usernamePattern = regexp.MustCompile(usernameAllowedChars)

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
	if !usernamePattern.MatchString(username) {
		return ErrInvalidUsernameChars
	}

	// Must start with a letter
	if !validationutil.ValidateStartLetter(username, []string{validationutil.IS_LETTER}) {
		return ErrUsernameMustStartWithLetter
	}

	// Must end with a letter or number
	if !validationutil.ValidateEndChar(username, []string{validationutil.IS_LETTER, validationutil.IS_NUMBER}) {
		return ErrUsernameMustEndWithLetterOrNumber
	}

	// No consecutive dots or underscores
	if !validationutil.ValidateNoConsecutive(username) {
		return ErrUsernameNoConsecutiveSpecialChars
	}

	// No adjacent dot/underscore combinations
	if !validationutil.ValidateNoAdjacentSpecialChars(username) {
		return ErrUsernameNoAdjacentSpecialChars
	}

	// Check reserved names
	if !validationutil.ValidateNotReservedName(username) {
		return ErrReservedUsername
	}

	return nil
}
