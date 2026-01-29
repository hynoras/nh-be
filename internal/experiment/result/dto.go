package result

import "time"

type ExperimentResultResponseDto struct {
	ID              string    `json:"id"`
	ExperimentID    string    `json:"experiment_id"`
	Outcome         string    `json:"outcome"`
	Summary         string    `json:"summary"`
	OutcomeReason   string    `json:"outcome_reason"`
	ConfidenceLevel string    `json:"confidence_level"`
	Version         int       `json:"version"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type CreateResultDto struct {
	ExperimentID    string `json:"experiment_id" binding:"required,uuid"`
	Outcome         string `json:"outcome" binding:"required,oneof=success failure inconclusive"`
	Summary         string `json:"summary" binding:"required,min=10"`
	OutcomeReason   string `json:"outcome_reason" binding:"required,min=10"`
	ConfidenceLevel string `json:"confidence_level" binding:"required,oneof=low medium high"`
}

type UpdateResultDto struct {
	Version         int    `json:"version" binding:"required,min=1"`
	Outcome         string `json:"outcome" binding:"omitempty,oneof=success failure inconclusive"`
	Summary         string `json:"summary" binding:"omitempty,min=10"`
	OutcomeReason   string `json:"outcome_reason" binding:"omitempty,min=10"`
	ConfidenceLevel string `json:"confidence_level" binding:"omitempty,oneof=low medium high"`
}
