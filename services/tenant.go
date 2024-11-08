package services

import (
	"fmt"
	"go-multitenant-boilerplate/models/tenant"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TenantService struct {
	db     *gorm.DB
	logger *zap.Logger
}

func InitializeTenantService(db *gorm.DB, log *zap.Logger) *TenantService {
	return &TenantService{db: db, logger: log}
}

func (ts *TenantService) Create(data *tenant.Tenant) error {
	var existingTenant tenant.Tenant
	if tx := ts.db.Where("organization = ?", data.Organization).First(&existingTenant); tx.Error != nil {
		if tx.Error != gorm.ErrRecordNotFound {
			return tx.Error
		}
	}
	if existingTenant.ID != "" {
		return fmt.Errorf("organization already exists: %s", existingTenant.Organization)
	}
	return ts.db.Create(data).Error
}
