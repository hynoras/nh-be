package auth

import "github.com/google/uuid"

type LoginDto struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponseDto struct {
	ID    uuid.UUID `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}