package stringutil

import (
	"strings"

	"github.com/google/uuid"
)

func ExtractUsernameFromEmail(email string) string {
	return email[:strings.Index(email, "@")]
}

func ConvertUUIDToString(uuid uuid.UUID) string {
	return uuid.String()
}
