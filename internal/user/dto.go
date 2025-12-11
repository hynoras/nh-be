package user

import "time"

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
}

type UserResponseDto struct {
	ID string `json:"id"`
	Username string `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	PermissionGroups []PermissionGroupResponseDto `json:"permission_groups"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserDto struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	Role string `json:"role"`
	Permissions []string `json:"permissions"`
}

type UpdateUserDto struct {
	Username string `json:"username,omitempty"`
	Email string `json:"email,omitempty"`
	Role string `json:"role,omitempty"`
	PermissionGroups []string `json:"permission_groups,omitempty"`
}

type DeleteUsersDto struct {
	IDs []string `json:"ids"`
}