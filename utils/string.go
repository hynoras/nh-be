package utils

import "github.com/google/uuid"

func ValidateUUIDString(uuidStr string) (error) {
	if err := uuid.Validate(uuidStr); err != nil {
		return err
	}
	return nil
}

func ValidateUUIDStrings(uuidStrs []string) (error) {
	for _, uuidToValidate := range uuidStrs {
		if err := uuid.Validate(uuidToValidate); err != nil {
			return err
		}
	}
	return nil
}

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

