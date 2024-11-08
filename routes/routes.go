package routes

import (
	"go-multitenant-boilerplate/controllers"
	"go-multitenant-boilerplate/middleware"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func StartRouter(logger *zap.Logger, controller *controllers.Controller, PORT string, GIN_MODE string) {
	gin.SetMode(GIN_MODE)
	router := gin.Default()
	router.Use()
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"status": false, "message": "Route not found"})
		c.Abort()
	})

	v1 := router.Group("/api/v1")
	{
		v1.POST("signup", controller.SignUp)
		v1.POST("login", controller.Login)
	}
	UserRoutes(v1, controller)

	if err := router.Run(PORT); err != nil {
		logger.Fatal("Failed to start server: " + err.Error())
	}
}

func UserRoutes(v1 *gin.RouterGroup, controller *controllers.Controller) {
	userRoutes := v1.Group("users")
	{
		userRoutes.POST("", middleware.TenantAuthMiddleware(), controller.AddUser)
		userRoutes.GET(":id", middleware.TenantAuthMiddleware(), controller.GetUser)
		userRoutes.GET("", middleware.TenantAuthMiddleware(), controller.ListUsers)
		userRoutes.PUT(":id", middleware.TenantAuthMiddleware(), controller.UpdateUser)
		userRoutes.DELETE(":id", middleware.TenantAuthMiddleware(), controller.DeleteUser)
	}
}
