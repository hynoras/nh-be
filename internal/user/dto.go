package user

import "time"

type UserResponseDto struct {
	ID string `json:"id"`
	Username string `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateUserDto struct {
	Username string `json:"username,omitempty"`
	Email string `json:"email,omitempty"`
	Role string `json:"role,omitempty"`
}