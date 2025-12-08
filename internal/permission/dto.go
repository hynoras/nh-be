package permission

import "time"

type PermissionResponseDto struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreatePermissionDto struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdatePermissionDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
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
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permission_ids"`
}

type UpdatePermissionGroupDto struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permission_ids"`
}

type AssignUserGroupDto struct {
	UserID            string `json:"user_id" binding:"required"`
	PermissionGroupID string `json:"permission_group_id" binding:"required"`
}
