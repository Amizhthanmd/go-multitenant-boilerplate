package controllers

import (
	"go-multitenant-boilerplate/services"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Controller struct {
	logger        *zap.Logger
	AdminDB       *gorm.DB
	TenantDB      *gorm.DB
	TenantService *services.TenantService
	userService   *services.UserService
}

func InitializeController(
	logger *zap.Logger,
	adminDB *gorm.DB,
	tenantDB *gorm.DB,
	tenantService *services.TenantService,
	userService *services.UserService,
) *Controller {
	return &Controller{
		logger:        logger,
		AdminDB:       adminDB,
		TenantDB:      tenantDB,
		TenantService: tenantService,
		userService:   userService,
	}
}
