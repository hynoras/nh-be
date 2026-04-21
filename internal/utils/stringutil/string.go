package stringutil

import (
	"strings"
)

func ExtractUsernameFromEmail(email string) string {
	return email[:strings.Index(email, "@")]
}

func StringPtr(s string) *string { return &s }
