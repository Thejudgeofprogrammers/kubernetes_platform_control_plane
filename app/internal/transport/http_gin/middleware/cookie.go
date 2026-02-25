package middleware

import "github.com/gin-gonic/gin"

func AddUserID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := "1"
		ctx.Set("user_id", userID)
		ctx.Set("role", "owner")
		ctx.Next()
	}
}