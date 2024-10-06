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
			ID: "add_users_table",
			Migrate: func(d *gorm.DB) error {
				return d.AutoMigrate(&tenant.User{})
			},
			Rollback: func(d *gorm.DB) error {
				return d.Migrator().DropTable(&tenant.User{})
			},
		},
	})

	if err := migrate.Migrate(); err != nil {
		log.Println("Failed to migrate")
	}
	return migrate
}
