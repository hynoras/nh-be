package observation

import (
	exp "nh-be/internal/features/experiment"
	proc "nh-be/internal/features/procedure"
	"time"

	"github.com/google/uuid"
)

type Observation struct {
	ID                    uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ObservedAt            time.Time    `gorm:"type:timestamp with time zone;not null;index"`
	Title                 string       `gorm:"type:text;not null"`
	Notes                 *string      `gorm:"type:text"`
	PreviousObservationID *Observation `gorm:"foreignKey:PreviousObservationID;references:ID;OnDelete:SET NULL"`
	CreatedBy             uuid.UUID    `gorm:"type:uuid;not null"`
	CreatedAt             time.Time    `gorm:"type:timestamp with time zone;not null;default:now()"`

	ExperimentID    uuid.UUID          `gorm:"type:uuid;not null;index"`
	Experiment      exp.Experiment     `gorm:"foreignKey:ExperimentID;references:ID;OnDelete:CASCADE"`
	ProcedureStepID *uuid.UUID         `gorm:"type:uuid;index"`
	ProcedureStep   proc.ProcedureStep `gorm:"foreignKey:ProcedureStepID;references:ID;OnDelete:SET NULL"`
}

//TODO: Uncomment when implementing measurement
// type ObservationMeasurement struct {
// 	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
// 	CreatedAt time.Time `gorm:"type:timestamp with time zone;not null;default:now()"`

// 	ValueFloat  *float64 `gorm:"type:numeric"`
// 	ValueInt    *int64   `gorm:"type:bigint"`
// 	ValueString *string  `gorm:"type:text"`
// 	ValueBool   *bool    `gorm:"type:boolean"`

// 	MeasurementParameterID uuid.UUID            `gorm:"not null;uniqueIndex:idx_obs_measure"`
// 	MeasurementParameter   MeasurementParameter `gorm:"foreignKey:MeasurementParameterID;references:ID;OnDelete:RESTRICT"`

// 	ObservationID uuid.UUID   `gorm:"not null;uniqueIndex:idx_obs_measure"`
// 	Observation   Observation `gorm:"foreignKey:ObservationID;references:ID;OnDelete:CASCADE"`
// }

//TODO: Move these two into measurement module
// type MeasurementParameter struct {
// 	ID        uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
// 	Name      string              `gorm:"type:varchar(50);not null;uniqueIndex:idx_obs_measure_type_def_by"`
// 	DataType  MeasurementDataType `gorm:"type:text;not null"`
// 	DefinedBy *uuid.UUID          `gorm:"type:uuid;not null;uniqueIndex:idx_obs_measure_type_def_by"`
// 	IsSystem  bool                `gorm:"not null;default:false"`

// 	DefaultUnitID *uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_obs_measure_type_def_by"`
// 	DefaultUnit   MeasurementUnit `gorm:"foreignKey:DefaultUnitID;references:ID;OnDelete:SET NULL"`
// }

// type MeasurementUnit struct {
// 	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
// 	Name      string     `gorm:"type:varchar(50);not null;uniqueIndex:idx_obs_measure_unit_def_by"`
// 	Symbol    string     `gorm:"type:varchar(10);not null;uniqueIndex:idx_obs_measure_unit_def_by"`
// 	DefinedBy *uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_obs_measure_unit_def_by"`
// 	IsSystem  bool       `gorm:"not null;default:false"`
// }
