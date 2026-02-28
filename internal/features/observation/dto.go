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

type ObservationsResponseDto struct {
	ID                    string  `json:"id"`
	ObservedAt            string  `json:"observed_at"`
	Title                 string  `json:"title"`
	Notes                 *string `json:"notes"`
	PreviousObservationID *string `json:"previous_observation_id"`
	CreatedBy             string  `json:"created_by"`
	CreatedAt             string  `json:"created_at"`
}
