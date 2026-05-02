package impl

import (
	"fmt"
	"time"

	"control_plane/internal/domain"
	"control_plane/internal/logger"
	JWTService "control_plane/internal/service/jwt"

	"github.com/golang-jwt/jwt/v5"
)

type jwtService struct {
	secret string
	exp    int
	log    logger.Logger
}

func NewJWTService(secret string, exp int, log logger.Logger) JWTService.JWTService {
	return &jwtService{
		secret: secret,
		exp:    exp,
		log:    log,
	}
}

func (s *jwtService) GenerateAccessToken(userID string, role domain.AccessRole) (string, error) {

	s.log.Debug("generate access token",
		"user_id", userID,
		"role", role,
	)

	claims := JWTService.Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.exp) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.secret))

	if err != nil {
		s.log.Error("failed to sign jwt",
			"user_id", userID,
			"role", role,
			"error", err,
		)
		return "", err
	}

	s.log.Debug("access token generated",
		"user_id", userID,
		"role", role,
		"exp", claims.ExpiresAt.Time,
	)

	return signed, nil
}

func (s *jwtService) Parse(tokenStr string) (*JWTService.Claims, error) {
	s.log.Debug("parse jwt started")
	token, err := jwt.ParseWithClaims(tokenStr, &JWTService.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})

	if err != nil {
		s.log.Warn("jwt parse failed",
			"error", err,
		)
		return nil, err
	}

	claims, ok := token.Claims.(*JWTService.Claims)
	if !ok {
		s.log.Error("invalid jwt claims type")
		return nil, fmt.Errorf("invalid claims type")
	}

	s.log.Debug("jwt parsed",
		"user_id", claims.UserID,
		"role", claims.Role,
		"exp", claims.ExpiresAt.Time,
	)

	return claims, nil
}
