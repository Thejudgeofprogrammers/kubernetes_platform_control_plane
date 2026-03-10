package jwt

import "control_plane/internal/domain"

type JWTService interface {
	GenerateAccessToken(userID string, role domain.AccessRole) (string, error)
	Parse(tokenStr string) (*Claims, error)
}
