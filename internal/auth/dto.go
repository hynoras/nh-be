package auth

import (
	"time"

	"github.com/google/uuid"
)

type LoginDto struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponseDto struct {
	ID uuid.UUID `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type LoginResponseDto struct {
	User UserResponseDto `json:"user"`
}

type ChangePasswordDto struct {
	NewPassword string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}