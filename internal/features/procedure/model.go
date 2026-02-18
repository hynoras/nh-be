package procedure

import (
	"nh-be/internal/features/experiment"
	"time"

	"github.com/google/uuid"
)

type Procedure struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title       string    `gorm:"type:text;not null;unique"`
	Description string    `gorm:"type:text"`

	Version  int        `gorm:"type:int;not null;default:1"`
	ParentID *uuid.UUID `gorm:"type:uuid"` // previous version

	CreatedAt time.Time  `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt *time.Time `gorm:"type:timestamp"`

	Steps       []ProcedureStep                 `gorm:"constraint:OnDelete:CASCADE"`
	Experiments []ProcedureExperimentAssignment `gorm:"constraint:OnDelete:CASCADE"`
}

type ProcedureStep struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProcedureID uuid.UUID `gorm:"type:uuid;not null;index"`
	Procedure   Procedure `gorm:"foreignKey:ProcedureID;references:ID;OnDelete:CASCADE"`

	Index       int    `gorm:"type:int;not null;default:0"`
	Title       string `gorm:"type:text;not null"`
	Description string `gorm:"type:text"`

	StepType string `gorm:"type:varchar(20);not null;index"`

	CreatedAt time.Time  `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt *time.Time `gorm:"type:timestamp"`
}

type ProcedureExperimentAssignment struct {
	ID           uuid.UUID             `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProcedureID  uuid.UUID             `gorm:"type:uuid;not null;uniqueIndex:idx_proc_exp"`
	ExperimentID uuid.UUID             `gorm:"type:uuid;not null;uniqueIndex:idx_proc_exp"`
	Procedure    Procedure             `gorm:"foreignKey:ProcedureID;references:ID;OnDelete:CASCADE"`
	Experiment   experiment.Experiment `gorm:"foreignKey:ExperimentID;references:ID;OnDelete:CASCADE"`

	CreatedAt time.Time  `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt *time.Time `gorm:"type:timestamp"`
}
