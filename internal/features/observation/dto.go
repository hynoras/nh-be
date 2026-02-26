package observation

import (
	"time"

	"github.com/google/uuid"
)

type ObservationMetadata struct {
	ID                    uuid.UUID
	ObservedAt            time.Time
	Title                 string
	Notes                 *string
	PreviousObservationID *uuid.UUID
	CreatedBy             uuid.UUID
	CreatedAt             time.Time
}
