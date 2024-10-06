package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeAdminDB() *gorm.DB {
	// Initialize the admin database
	dsn := fmt.Sprintf("%s%s", os.Getenv("POSTGRES_DB_URL"), os.Getenv("ADMIN_DB"))
	adminDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Error initializing admin database")
	}
	return adminDB
}

func InitializeTenantDB() *gorm.DB {
	// Initialize the tenant database
	dsn := fmt.Sprintf("%s%s", os.Getenv("POSTGRES_DB_URL"), os.Getenv("TENANT_DB"))
	tenantDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Error initializing tenant database")
	}
	return tenantDB
}

func CreateSchema(db *gorm.DB, org string) error {
	err := db.Exec(fmt.Sprintf(`CREATE SCHEMA "%s";`, org)).Error
	if err != nil {
		log.Println("Error creating schema:", err)
		return err
	}
	return nil
}
