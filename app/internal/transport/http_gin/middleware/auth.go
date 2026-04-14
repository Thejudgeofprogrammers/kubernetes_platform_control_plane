package middleware

import (
	jwt "control_plane/internal/service/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(jwtService jwt.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "unauthorize"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claim, err := jwtService.Parse(tokenStr)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		ctx.Set("user_id", string(claim.UserID))
		ctx.Set("role", string(claim.Role))

		ctx.Next()
	}
}
