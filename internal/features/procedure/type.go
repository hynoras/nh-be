package procedure

import "github.com/google/uuid"

type StepMetadata struct {
	ID      uuid.UUID
	Version int
}
