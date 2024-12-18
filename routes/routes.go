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
	RoleRoutes(v1, controller)
	PermissionRoutes(v1, controller)

	if err := router.Run(PORT); err != nil {
		logger.Fatal("Failed to start server: " + err.Error())
	}
}

func UserRoutes(v1 *gin.RouterGroup, controller *controllers.Controller) {
	userRoutes := v1.Group("users")
	{
		userRoutes.POST("", middleware.AuthMiddleware("users:write"), controller.AddUser)
		userRoutes.GET(":id", middleware.AuthMiddleware("users:read"), controller.GetUser)
		userRoutes.GET("", middleware.AuthMiddleware("users:list"), controller.ListUsers)
		userRoutes.PUT(":id", middleware.AuthMiddleware("users:update"), controller.UpdateUser)
		userRoutes.DELETE(":id", middleware.AuthMiddleware("users:delete"), controller.DeleteUser)
	}
}

func RoleRoutes(v1 *gin.RouterGroup, controller *controllers.Controller) {
	userRoutes := v1.Group("roles")
	{
		userRoutes.POST("", middleware.AuthMiddleware(""), controller.AddRoles)
		userRoutes.GET("", middleware.AuthMiddleware(""), controller.ListRoles)
		userRoutes.PUT(":id", middleware.AuthMiddleware(""), controller.UpdateRoles)
		userRoutes.DELETE(":id", middleware.AuthMiddleware(""), controller.DeleteRoles)
	}
}

func PermissionRoutes(v1 *gin.RouterGroup, controller *controllers.Controller) {
	userRoutes := v1.Group("permissions")
	{
		userRoutes.GET("", middleware.AuthMiddleware(""), controller.ListPermissions)
	}
}
