package auth

import (
	"nh-be/internal/user"
	"time"

	"github.com/google/uuid"
)

type VerificationTokenType string

const (
	VerifyEmail   VerificationTokenType = "verify_email"
	ResetPassword VerificationTokenType = "reset_password"
)

type VerificationToken struct {
	ID        uint                  `gorm:"primaryKey;autoIncrement"`
	UserID    uuid.UUID             `gorm:"type:uuid;not null;uniqueIndex"`
	User      user.User             `gorm:"constraint:OnDelete:CASCADE"`
	Type      VerificationTokenType `gorm:"type:varchar(20);not null;index"`
	CodeHash  string                `gorm:"type:char(64);not null;uniqueIndex"`
	ExpireAt  time.Time             `gorm:"not null;index"`
	CreatedAt time.Time             `gorm:"autoCreateTime"`
}
