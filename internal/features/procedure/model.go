package procedure

import (
	"time"

	"github.com/google/uuid"
)

type Procedure struct {
	ID          uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Title       string          `gorm:"type:text;not null"`
	Description *string         `gorm:"type:text"`
	ParentID    *uuid.UUID      `gorm:"type:uuid"` // previous version
	CreatedAt   time.Time       `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt   *time.Time      `gorm:"type:timestamp"`
	Version     int             `gorm:"type:int;not null;default:1"`
	Steps       []ProcedureStep `gorm:"constraint:OnDelete:CASCADE"`
}

type ProcedureStep struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ProcedureID uuid.UUID  `gorm:"type:uuid;not null;index"`
	Procedure   Procedure  `gorm:"foreignKey:ProcedureID;references:ID;OnDelete:CASCADE"`
	Index       int        `gorm:"type:int;not null"`
	Title       string     `gorm:"type:text;not null"`
	Description *string    `gorm:"type:text"`
	IsOptional  bool       `gorm:"type:boolean;not null;default:false"`
	StepType    StepType   `gorm:"type:varchar(20);not null;index"`
	WaitTime    *int       `gorm:"type:int"` //in milliseconds
	CreatedAt   time.Time  `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt   *time.Time `gorm:"type:timestamp"`
	Version     int        `gorm:"type:int;not null;default:1"`
}
