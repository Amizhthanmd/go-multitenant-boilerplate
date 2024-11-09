package middleware

import (
	"go-multitenant-boilerplate/helpers"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		organization := ctx.GetHeader("Organization")
		if organization == "" {
			ctx.JSON(401, gin.H{"error": "Organization header missing or invalid"})
			ctx.Abort()
			return
		}
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(401, gin.H{"error": "Authorization header missing or invalid"})
			ctx.Abort()
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := ParseJWT(tokenString)
		if err != nil {
			ctx.JSON(401, gin.H{"status": false, "message": "Invalid token"})
			ctx.Abort()
			return
		}

		if err := claims.Valid(); err != nil {
			ctx.JSON(401, gin.H{"status": false, "message": "Token expired or invalid"})
			ctx.Abort()
			return
		}

		if permission != "" {
			if !helpers.SliceContains(claims.Permissions, permission) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Permission is not allowed"})
				ctx.Abort()
				return
			}
		}

		ctx.Set("user_id", claims.ID)
		ctx.Set("email", claims.Email)
		ctx.Set("organization", claims.Organization)
		ctx.Next()
	}
}
