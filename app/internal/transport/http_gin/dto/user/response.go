package dto

import (
	"control_plane/internal/domain"
)

type UserResponse struct {
	ID        string            `json:"id"`
	Email     string            `json:"email"`
	FullName  string            `json:"full_name"`
	Role      domain.AccessRole `json:"role"`
	CreatedAt string            `json:"created_at"`
}
