package main

import (
	"fmt"
	"go-multitenant-boilerplate/db/migrations"
	"go-multitenant-boilerplate/helpers"
	"go-multitenant-boilerplate/models/tenant"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Choose a migration DB to migrate :\n", "1 : Admin\n", "2 : Tenant")
	var dbChoice int
	fmt.Scan(&dbChoice)

	helpers.LoadEnv()

	adminDbName := os.Getenv("ADMIN_DB")
	tenantDbName := os.Getenv("TENANT_DB")
	initialDB, err := gorm.Open(postgres.Open(os.Getenv("POSTGRES_DB_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to PostgreSQL server")
	}

	CheckAndCreateDatabase(initialDB, adminDbName)
	dsn := fmt.Sprintf("%s%s", os.Getenv("POSTGRES_DB_URL"), os.Getenv("ADMIN_DB"))
	adminDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	if dbChoice == 1 {
		if err := migrations.MigrateAdminDB(adminDb).Migrate(); err != nil {
			log.Println("Failed to migrate admin database")
		}
		fmt.Println("Admin database migrated successfully")
	} else if dbChoice == 2 {
		CheckAndCreateDatabase(initialDB, tenantDbName)
		var tenant []tenant.Tenant
		if tx := adminDb.Table("tenants").Find(&tenant); tx.Error != nil {
			log.Println("Error: ", tx.Error)
		}

		dsn := fmt.Sprintf("%s%s", os.Getenv("POSTGRES_DB_URL"), os.Getenv("TENANT_DB"))
		tenantDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect tenant database")
		}

		for _, t := range tenant {
			tenantDb.Exec(fmt.Sprintf(`set search_path="%s"`, t.Organization))
			if err := migrations.MigrateTenantDB(tenantDb).Migrate(); err != nil {
				log.Println("Failed to migrate tenant database for tenant ID: ", t.ID)
			}
			fmt.Println("Tenant database migrated successfully for tenant ID: ", t.ID)
		}
	}
}

func CheckAndCreateDatabase(initialDB *gorm.DB, dbName string) {
	var exists bool
	err := initialDB.Raw("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = ?)", dbName).Scan(&exists).Error
	if err != nil {
		log.Fatal("failed to check if database exists:", err)
	}

	if !exists {
		err = initialDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error
		if err != nil {
			log.Fatal("failed to create database:", err)
		}
		log.Println("Database created successfully:", dbName)
	} else {
		log.Println("Database already exists:", dbName)
	}
}
