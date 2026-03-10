package jwt

import (
	"control_plane/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string
	Role   domain.AccessRole
	jwt.RegisteredClaims
}
