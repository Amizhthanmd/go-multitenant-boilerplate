package migrations

import (
	"go-multitenant-boilerplate/models/tenant"
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func MigrateTenantDB(db *gorm.DB) *gormigrate.Gormigrate {
	migrate := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "001_add_users_table",
			Migrate: func(d *gorm.DB) error {
				return d.AutoMigrate(&tenant.User{})
			},
			Rollback: func(d *gorm.DB) error {
				return d.Migrator().DropTable(&tenant.User{})
			},
		},
		{
			ID: "002_add_roles_table",
			Migrate: func(d *gorm.DB) error {
				return d.AutoMigrate(&tenant.Role{})
			},
			Rollback: func(d *gorm.DB) error {
				return d.Migrator().DropTable(&tenant.Role{})
			},
		},
		{
			ID: "003_add_permissions_table",
			Migrate: func(d *gorm.DB) error {
				return d.AutoMigrate(&tenant.Permission{})
			},
			Rollback: func(d *gorm.DB) error {
				return d.Migrator().DropTable(&tenant.Permission{})
			},
		},
	})

	if err := migrate.Migrate(); err != nil {
		log.Println("Failed to migrate")
	}
	return migrate
}

func MigratePermission(tenantDb *gorm.DB) error {
	permission := []tenant.Permission{
		{Name: "users:write"},
		{Name: "users:read"},
		{Name: "users:list"},
		{Name: "users:update"},
		{Name: "users:delete"},
	}
	return tenantDb.Create(&permission).Error
}
