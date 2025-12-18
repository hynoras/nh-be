package permission

import (
	"time"

	"github.com/google/uuid"
)

type PermissionResponseDto struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

type PermissionGroupResponseDto struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Permissions []PermissionResponseDto `json:"permissions"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

type CreatePermissionGroupDto struct {
	Name          string   `json:"name" binding:"required,min=5,max=50"`
	Description   string   `json:"description" binding:"omitempty"`
	Permissions []string `json:"permissions" binding:"required,dive"`
}

type UpdatePermissionGroupDto struct {
	Name          string   `json:"name" binding:"omitempty,min=5,max=50"`
	Description   string   `json:"description" binding:"omitempty"`
	Permissions []string `json:"permissions" binding:"required,dive,uuid4"`
}

type PermissionGroupInput struct {
	Name          string   
	Description   string   
	Permissions []uuid.UUID 
}