package main

import (
	"fmt"
	"go-multitenant-boilerplate/appInit"
	"go-multitenant-boilerplate/controllers"
	"go-multitenant-boilerplate/db"
	"go-multitenant-boilerplate/helpers"
	"go-multitenant-boilerplate/routes"
	"go-multitenant-boilerplate/services"
	"os"
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

	// Initialize Routes
	routes.StartRouter(logger, controller, PORT, GIN_MODE)
}
