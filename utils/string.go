package utils

import "github.com/google/uuid"

func ParseStringToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func ParseUUIDToString(id uuid.UUID) string {
	return id.String()
}