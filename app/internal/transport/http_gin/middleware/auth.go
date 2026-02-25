package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(enabled string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.ToLower(enabled) == "false" {
			ctx.Next()
			return
		}

		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.AbortWithStatusJSON(401, gin.H{
				"error": "unauthorized",
			})
			return
		}

		ctx.Next()
	}
}
