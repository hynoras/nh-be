package user

import (
	"nh-be/internal/permission"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                       uuid.UUID                    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username                 string                       `gorm:"type:varchar(255);unique;not null"`
	Email                    string                       `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password                 string                       `gorm:"type:varchar(255);not null"`
	AssignedPermissionGroups []permission.PermissionGroup `gorm:"many2many:user_permissions;joinForeignKey:UserID;joinReferences:PermissionGroupID;constraint:OnDelete:CASCADE"`
	IsVerified               bool                         `gorm:"type:boolean;not null;default:false"`
	CreatedAt                time.Time                    `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt                time.Time                    `gorm:"type:timestamp;not null;default:now()"`
}

type UserPermission struct {
	UserID            uuid.UUID                  `gorm:"primaryKey;type:uuid"`
	PermissionGroupID uuid.UUID                  `gorm:"primaryKey;type:uuid"`
	User              User                       `gorm:"constraint:OnDelete:CASCADE"`
	PermissionGroup   permission.PermissionGroup `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt         time.Time                  `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt         time.Time                  `gorm:"type:timestamp;not null;default:now()"`
}
