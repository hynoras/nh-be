package experiment

import (
	"time"

	proc "nh-be/internal/features/procedure"
	"nh-be/internal/features/user"

	"github.com/google/uuid"
)

type Experiment struct {
	ID uuid.UUID `gorm:"type:uuid;primaryKey"`

	// Human-readable identifier
	Identifier string `gorm:"type:varchar(50);uniqueIndex;not null"`

	// Metadata
	Title     string `gorm:"type:varchar(200);not null;uniqueIndex:idx_user_experiment_title"` //composite key with CreatedByID
	Objective string `gorm:"type:text"`

	// Workflow
	Status ExperimentStatus `gorm:"type:varchar(20);not null;index"`
	Type   ExperimentType   `gorm:"type:varchar(20);not null;index"`

	// Optimistic locking
	Version int `gorm:"not null;default:1"`

	// Ownership
	CreatedByID uuid.UUID `gorm:"not null;uniqueIndex:idx_user_experiment_title"` //composite key with Title
	CreatedBy   user.User

	UpdatedByID *uuid.UUID `gorm:"type:uuid;index"`
	UpdatedBy   *user.User

	// Scheduling
	PlannedStartAt *time.Time
	PlannedEndAt   *time.Time

	// Execution
	StartedAt   *time.Time
	CompletedAt *time.Time

	// Procedure template
	ProcedureID *uuid.UUID `gorm:"type:uuid;index"`
	Procedure   proc.Procedure

	CreatedAt time.Time
	UpdatedAt time.Time
}
