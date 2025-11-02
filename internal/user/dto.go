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