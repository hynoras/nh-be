package auth

import (
	"time"

	"github.com/google/uuid"
)

type LoginDto struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponseDto struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type LoginResponseDto struct {
	User UserResponseDto `json:"user"`
}

type ChangePasswordDto struct {
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type SignUpDto struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreatedTokenDto struct {
	UserID    uuid.UUID             `json:"user_id"`
	Type      VerificationTokenType `json:"type"`
	Token     string                `json:"token"`
	CreatedAt time.Time             `json:"created_at"`
	ExpireAt  time.Time             `json:"expire_at"`
}

type CreateVerificationTokenDto struct {
	Email string `json:"email" binding:"required,email"`
	Type  string `json:"type" binding:"required, oneof=verify_email reset_password"`
}
