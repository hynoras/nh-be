package permission

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name        string    `gorm:"type:text;not null;unique"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt   time.Time `gorm:"type:timestamp;not null;default:now()"`
}

type PermissionGroup struct {
	ID            uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name          string       `gorm:"type:text;not null;unique"`
	Description   string       `gorm:"type:text"`
	Permissions   []Permission `gorm:"many2many:permission_group_items;constraint:OnDelete:CASCADE"`
	CreatedAt     time.Time    `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt     time.Time    `gorm:"type:timestamp;not null;default:now()"`
}

