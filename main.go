package main

import (
	"fmt"
	"go-multitenant-boilerplate/appInit"
	"go-multitenant-boilerplate/controllers"
	"go-multitenant-boilerplate/db"
	"go-multitenant-boilerplate/helpers"
	"go-multitenant-boilerplate/middleware"
	"go-multitenant-boilerplate/services"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load ENV
	helpers.LoadEnv()

	var (
		PORT     = fmt.Sprintf(":%s", os.Getenv("PORT"))
		GIN_MODE = os.Getenv("GIN_MODE")
	)

	// Initialize logger
	logger := appInit.ZapLogger(GIN_MODE)
	defer logger.Sync()

	adminDB := db.InitializeAdminDB()
	tenantDB := db.InitializeTenantDB()
	tenantService := services.InitializeTenantService(adminDB, logger)
	userService := services.InitializeUserService(tenantDB, logger)
	// Controller
	controller := controllers.InitializeController(logger, adminDB, tenantDB, tenantService, userService)

	// Initialize Gin engine and middleware
	gin.SetMode(GIN_MODE)
	router := gin.Default()
	router.Use()
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"status": false, "message": "Route not found"})
		c.Abort()
	})

	logger.Info("Starting server on port " + PORT)

	v1 := router.Group("/api/v1")
	{
		v1.POST("signup", controller.SignUp)
		v1.POST("login", controller.Login)
	}

	userRoutes := v1.Group("users")
	{
		userRoutes.POST("", middleware.TenantAuthMiddleware(), controller.AddUser)
		userRoutes.GET("", middleware.TenantAuthMiddleware(), controller.GetUsers)
		userRoutes.PUT(":id", middleware.TenantAuthMiddleware(), controller.UpdateUser)
		userRoutes.DELETE(":id", middleware.TenantAuthMiddleware(), controller.DeleteUser)
	}

	if err := router.Run(PORT); err != nil {
		logger.Fatal("Failed to start server: " + err.Error())
	}
}
