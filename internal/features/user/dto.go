package user

import (
	"nh-be/internal/features/permission"
	"time"

	"github.com/google/uuid"
)

type UserResponseDto struct {
	ID               string                                  `json:"id"`
	Username         string                                  `json:"username"`
	Email            string                                  `json:"email"`
	PermissionGroups []permission.PermissionGroupResponseDto `json:"permission_groups"`
	CreatedAt        time.Time                               `json:"created_at"`
}

type MeResponseDto struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	Permissions []string  `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateUserDto struct {
	Username    string   `json:"username" binding:"min=3,max=30"`
	Email       string   `json:"email" binding:"required,email"`
	Password    string   `json:"password" binding:"required,min=8"`
	Permissions []string `json:"permissions" binding:"omitempty,dive,uuid"`
}

type UpdateUserDto struct {
	Username    string   `json:"username,omitempty" binding:"omitempty,min=3,max=30"`
	Email       string   `json:"email,omitempty" binding:"omitempty,email"`
	Permissions []string `json:"permissions,omitempty" binding:"omitempty,dive,uuid"`
}

type UserInput struct {
	Username    string
	Email       string
	Password    string
	Permissions []uuid.UUID
}

type DeleteUsersDto struct {
	IDs []string `json:"ids"`
}

type CreatedUserDto struct {
	ID         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}
