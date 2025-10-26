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
	Role string `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type LoginResponseDto struct {
	User UserResponseDto `json:"user"`
}