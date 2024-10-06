package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func TenantAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(401, gin.H{"status": false, "message": "Missing Authorization header"})
			ctx.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer")
		claims, err := ParseJWT(tokenStr)
		if err != nil {
			ctx.JSON(401, gin.H{"status": false, "message": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.ID)
		ctx.Set("email", claims.Email)
		ctx.Set("organization", claims.Organization)
		ctx.Next()
	}
}
