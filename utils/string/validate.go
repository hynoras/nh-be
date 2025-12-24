package string

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

const (
	IS_LOWERCASE = "is_lowercase"
	IS_UPPERCASE = "is_uppercase"
	IS_NUMBER    = "is_number"
	IS_LETTER    = "is_letter"
)

func ValidateWithOptions(stringToValidate rune, options []string) bool {
	for _, option := range options {
		switch option {
		case IS_LOWERCASE:
			if unicode.IsLower(stringToValidate) {
				return true
			}
		case IS_UPPERCASE:
			if unicode.IsUpper(stringToValidate) {
				return true
			}
		case IS_NUMBER:
			if unicode.IsDigit(stringToValidate) {
				return true
			}
		case IS_LETTER:
			if unicode.IsLetter(stringToValidate) {
				return true
			}
		}
	}
	return false
}

/*
ValidateStartLetter validates the start letter of a string
Checks if the first character matches any of the provided options (OR logic)
@param stringToValidate string - the string to validate
@param options []string - list of validation options (IS_LOWERCASE, IS_UPPERCASE, IS_NUMBER, IS_LETTER)
@return bool - true if the first character matches any option, false otherwise
*/
func ValidateStartLetter(stringToValidate string, options []string) bool {
	if len(stringToValidate) == 0 {
		return false
	}
	if len(options) == 0 {
		return false
	}

	firstRune := []rune(stringToValidate)[0]
	return ValidateWithOptions(firstRune, options)
}

func ValidateUUIDString(uuidStr string) error {
	if err := uuid.Validate(uuidStr); err != nil {
		return err
	}
	return nil
}

/*
ValidateAllowedChar validates the allowed characters of a string
@param stringToValidate string
@param allowedCharRegex string
@return bool
*/
func ValidateAllowedChar(stringToValidate string, allowedCharRegex string) bool {
	return regexp.MustCompile(allowedCharRegex).MatchString(stringToValidate)
}

func ValidateNoConsecutive(stringToValidate string) bool {
	return !strings.Contains(stringToValidate, "..") && !strings.Contains(stringToValidate, "__")
}

func ValidateEndChar(stringToValidate string, options []string) bool {
	if len(stringToValidate) == 0 {
		return false
	}
	lastRune := []rune(stringToValidate)[len(stringToValidate)-1]
	return ValidateWithOptions(lastRune, options)
}

func ValidateUUIDStrings(uuidStrs []string) error {
	for _, uuidToValidate := range uuidStrs {
		if err := uuid.Validate(uuidToValidate); err != nil {
			return err
		}
	}
	return nil
}

// ValidateNoAdjacentSpecialChars checks if a string contains adjacent special characters (dot and underscore)
// Returns true if no adjacent special chars are found, false otherwise
func ValidateNoAdjacentSpecialChars(s string) bool {
	return !strings.Contains(s, "._") && !strings.Contains(s, "_.")
}

// ReservedNames contains names that cannot be used
var ReservedNames = []string{
	"admin",
	"root",
	"system",
	"support",
	"api",
	"null",
	"undefined",
	"superadmin",
}

// ValidateNotReservedName checks if a name is not in the reserved names list
// Returns true if the name is not reserved, false otherwise
func ValidateNotReservedName(name string) bool {
	lowerName := strings.ToLower(name)
	for _, reserved := range ReservedNames {
		if lowerName == reserved {
			return false
		}
	}
	return true
}
