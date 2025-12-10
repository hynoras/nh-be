package utils

import "github.com/google/uuid"

func ParseStringToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func ParseStringsToUUIDs(ss []string) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	for _, s := range ss {
		id, err := uuid.Parse(s)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func ParseUUIDToString(id uuid.UUID) string {
	return id.String()
}

