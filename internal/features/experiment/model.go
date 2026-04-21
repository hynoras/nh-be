package experiment

import (
	"time"

	proc "nh-be/internal/features/procedure"
	"nh-be/internal/features/user"

	"github.com/google/uuid"
)

type Experiment struct {
	ID        uuid.UUID        `gorm:"type:uuid;primaryKey"`
	Title     string           `gorm:"not null"`
	Objective string           `gorm:"type:text"`
	Status    ExperimentStatus `gorm:"type:varchar(20);not null;index;default:draft"`
	Type      ExperimentType   `gorm:"type:varchar(20);not null;index"`
	Version   int              `gorm:"not null;default:1"`

	CreatedByID uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedBy   user.User `gorm:"foreignKey:CreatedByID"`

	StartedAt   *time.Time
	CompletedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time

	ProcedureID *uuid.UUID     `gorm:"type:uuid;index"`
	Procedure   proc.Procedure `gorm:"foreignKey:ProcedureID"`

	// Relations
	// Steps        []ExperimentStep        `gorm:"constraint:OnDelete:CASCADE"`
	// Materials    []ExperimentMaterial    `gorm:"constraint:OnDelete:CASCADE"`
	// Conditions   []ExperimentCondition   `gorm:"constraint:OnDelete:CASCADE"`
	// Observations []ExperimentObservation `gorm:"constraint:OnDelete:CASCADE"`
	// Results      []ExperimentResult      `gorm:"constraint:OnDelete:CASCADE"`
	// Notes        []ExperimentNote        `gorm:"constraint:OnDelete:CASCADE"`
}
