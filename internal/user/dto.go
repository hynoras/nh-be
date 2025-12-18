package user

import (
	"time"

	"github.com/google/uuid"
)

type PermissionResponseDto struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PermissionGroupResponseDto struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Permissions []PermissionResponseDto `json:"permissions"`
}

type UserResponseDto struct {
	ID               string                       `json:"id"`
	Username         string                       `json:"username"`
	Email            string                       `json:"email"`
	Role             string                       `json:"role"`
	PermissionGroups []PermissionGroupResponseDto `json:"permission_groups"`
	CreatedAt        time.Time                    `json:"created_at"`
}

type CreateUserDto struct {
	Username    string   `json:"username" binding:"min=3,max=20,lowercase"`
	Email       string   `json:"email" binding:"required,email,unique"`
	Password    string   `json:"password" binding:"required,min=8"`
	Role        string   `json:"role" binding:"omitempty,oneof=admin user"`
	Permissions []string `json:"permissions" binding:"omitempty,dive,uuid"`
}

type UpdateUserDto struct {
	Username    string   `json:"username,omitempty" binding:"omitempty,min=3,max=20,lowercase"`
	Email       string   `json:"email,omitempty" binding:"omitempty,email"`
	Role        string   `json:"role,omitempty" binding:"omitempty,oneof=admin user"`
	Permissions []string `json:"permissions,omitempty" binding:"omitempty,dive,uuid"`
}

type UserInput struct {
	Username    string
	Email       string
	Password    string
	Role        string
	Permissions []uuid.UUID
}

type DeleteUsersDto struct {
	IDs []string `json:"ids"`
}
