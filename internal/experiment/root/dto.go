package root

import "time"

type ExperimentsResponseDto struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Objective string    `json:"objective"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type ExperimentResponseDto struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Objective   string     `json:"objective"`
	Type        string     `json:"type"`
	Status      string     `json:"status"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	StartedAt   *time.Time `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"`
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
	Status string `json:"status" binding:"required"`
}
