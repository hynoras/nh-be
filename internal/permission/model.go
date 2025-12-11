package permission

import (
	"nh-be/internal/user"
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
	AssignedUsers []user.User  `gorm:"many2many:user_permissions;joinForeignKey:PermissionGroupID;joinReferences:UserID;constraint:OnDelete:CASCADE"`
	CreatedAt     time.Time    `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt     time.Time    `gorm:"type:timestamp;not null;default:now()"`
}

type UserPermission struct {
	UserID            uuid.UUID `gorm:"primaryKey;type:uuid"`
	PermissionGroupID uuid.UUID `gorm:"primaryKey;type:uuid"`
	User             user.User        `gorm:"constraint:OnDelete:CASCADE"`
    PermissionGroup  PermissionGroup  `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt         time.Time `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt         time.Time `gorm:"type:timestamp;not null;default:now()"`
}
