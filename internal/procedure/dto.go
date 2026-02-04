package procedure

import (
	"time"
)

type UsedByExperiment struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Objective string `json:"objective"`
}

type Steps struct {
	ID          string `json:"step_id"`
	ProcedureID string `json:"procedure_id"`

	Index       int    `json:"step_order"`
	Title       string `json:"title"`
	Description string `json:"description"`

	StepType string `json:"type"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProcedureListResponseDto struct {
	ID                string             `json:"id"`
	Title             string             `json:"title"`
	Description       string             `json:"description"`
	UsedByExperiments []UsedByExperiment `json:"used_by_experiments"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

type ProcedureResponseDto struct {
	ID                string             `json:"id"`
	Title             string             `json:"title"`
	Description       string             `json:"description"`
	UsedByExperiments []UsedByExperiment `json:"used_by_experiments"`
	Steps             []Steps            `json:"steps"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}
