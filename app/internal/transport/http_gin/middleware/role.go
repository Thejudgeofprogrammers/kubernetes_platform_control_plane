package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(enabled string, allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.ToLower(enabled) == "false" {
			ctx.Next()
			return
		}

		role := ctx.GetString("role")
		for _, r := range allowedRoles {
			if role == r {
				ctx.Next()
				return
			}
		}

		ctx.AbortWithStatusJSON(403, gin.H{
			"error": "forbidden",
		})
	}
}
