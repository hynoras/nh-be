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

type CreatedObservationResponseDto struct {
	ID                    string  `json:"id"`
	ExperimentID          string  `json:"experiment_id"`
	ProcedureStepID       *string `json:"procedure_step_id"`
	ObservedAt            string  `json:"observed_at"`
	Title                 string  `json:"title"`
	Notes                 *string `json:"notes"`
	PreviousObservationID *string `json:"previous_observation_id"`
	CreatedBy             string  `json:"created_by"`
	CreatedAt             string  `json:"created_at"`
}

type CreateObservationDto struct {
	Title                 string    `json:"title" binding:"required,max=150"`
	Notes                 *string   `json:"notes" binding:"omitempty,max=5000"`
	PreviousObservationID *string   `json:"previous_observation_id" binding:"omitempty,uuidv4"`
	ObservedAt            time.Time `json:"observed_at" binding:"required"`
}

type CreateObservationInput struct {
	ObservedAt            time.Time  `json:"observed_at"`
	Title                 string     `json:"title"`
	Notes                 *string    `json:"notes"`
	PreviousObservationID *uuid.UUID `json:"previous_observation_id"`
}
