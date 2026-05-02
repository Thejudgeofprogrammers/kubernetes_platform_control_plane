package mapper

import (
	userDTO "control_plane/internal/transport/http_gin/dto/user"
	"control_plane/internal/domain"
	"time"
)

func ToUserResponse(u domain.User) userDTO.UserResponse {
	return userDTO.UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}
