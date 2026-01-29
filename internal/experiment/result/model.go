package result

import (
	"time"

	"github.com/google/uuid"
)

type Outcome string
type ConfidenceLevel string

const (
	OutcomeSuccess      Outcome = "success"
	OutcomeFailure      Outcome = "failure"
	OutcomeInconclusive Outcome = "inconclusive"
)

const (
	ConfidenceLow    ConfidenceLevel = "low"
	ConfidenceMedium ConfidenceLevel = "medium"
	ConfidenceHigh   ConfidenceLevel = "high"
)

type ExperimentResult struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	ExperimentID uuid.UUID `gorm:"type:uuid;not null;unique;index"`

	Outcome         Outcome         `gorm:"type:varchar(20);not null"`
	Summary         string          `gorm:"type:text;not null"`
	OutcomeReason   string          `gorm:"type:text;not null"`
	ConfidenceLevel ConfidenceLevel `gorm:"type:varchar(20);not null"`
	Version         int             `gorm:"not null;default:1"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
