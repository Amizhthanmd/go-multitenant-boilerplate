package migrations

import (
	"go-multitenant-boilerplate/models/tenant"
	"log"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func MigrateAdminDB(db *gorm.DB) *gormigrate.Gormigrate {
	migrate := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "001_create_tenants_table",
			Migrate: func(d *gorm.DB) error {
				return d.AutoMigrate(&tenant.Tenant{})
			},
			Rollback: func(d *gorm.DB) error {
				return d.Migrator().DropTable(&tenant.Tenant{})
			},
		},
	})

	if err := migrate.Migrate(); err != nil {
		log.Println("Failed to migrate")
	}
	return migrate
}
