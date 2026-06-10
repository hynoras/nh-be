package experiment

import (
	"time"
)

type ExperimentsResponseDto struct {
	Identifier    string    `json:"identifier"`
	Title         string    `json:"title"`
	Objective     string    `json:"objective"`
	Status        string    `json:"status"`
	Type          string    `json:"type"`
	Creator       string    `json:"created_by"`
	Updater       string    `json:"updated_by"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	ProcedureName string    `json:"procedure_name"`
}

type ExperimentResponseDto struct {
	ID             string     `json:"id"`
	Identifier     string     `json:"identifier"`
	Title          string     `json:"title"`
	Objective      string     `json:"objective"`
	Type           string     `json:"type"`
	Status         string     `json:"status"`
	CreatedByID    string     `json:"created_by_id"`
	CreatedBy      string     `json:"created_by"`
	UpdatedByID    string     `json:"updated_by_id"`
	UpdatedBy      string     `json:"updated_by"`
	PlannedStartAt *time.Time `json:"planned_start_at"`
	PlannedEndAt   *time.Time `json:"planned_end_at"`
	StartedAt      *time.Time `json:"started_at"`
	CompletedAt    *time.Time `json:"completed_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	Version        int        `json:"version"`
}

type ExperimentsQueryDto struct {
	Identifier    string
	Title         string
	Objective     string
	Status        ExperimentStatus
	Type          ExperimentType
	Creator       string
	Updater       string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ProcedureName string
}

type CreateExperimentDto struct {
	Title     string `json:"title" binding:"required,min=3,max=200"`
	Type      string `json:"type" binding:"required"`
	Objective string `json:"objective" binding:"required,min=5,max=255"`
}

type UpdateExperimentDto struct {
	Title     string `json:"title" binding:"omitempty,min=3,max=200"`
	Type      string `json:"type" binding:"omitempty"`
	Objective string `json:"objective" binding:"omitempty,min=5,max=255"`
}

type UpdateExperimentStatusDto struct {
	Status  string `json:"status" binding:"required"`
	Version int    `json:"version" binding:"required"`
}

type AssignProcedureToExperimentDto struct {
	Version int `json:"version" binding:"required"`
}
