package procedure

import (
	"time"
)

type UsedByExperiment struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Objective string `json:"objective"`
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
	Steps             []ProcedureStep    `json:"steps"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}
