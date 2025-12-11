package user

import (
	"time"

	"github.com/google/uuid"
)

//create another permission to avoid import cycle
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

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Username  string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string    `gorm:"type:varchar(255);not null"`
	Role      	string    `gorm:"type:varchar(255);not null;default:user"`
	AssignedPermissionGroups []PermissionGroup `gorm:"many2many:user_permissions;joinForeignKey:UserID;joinReferences:PermissionGroupID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `gorm:"type:timestamp;not null;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamp;not null;default:now()"`
}