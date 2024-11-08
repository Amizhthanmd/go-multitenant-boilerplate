package main

import (
	"fmt"
	"go-multitenant-boilerplate/db/migrations"
	"go-multitenant-boilerplate/models/tenant"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load("../../../.env")
	if err != nil {
		fmt.Println("Failed to load .env :", err)
	}
	fmt.Print("1 : Create Database\n", "2 : Run Migration\n", "Choose the option : ")
	var DbOption int
	fmt.Scan(&DbOption)
	dbNames := []string{os.Getenv("ADMIN_DB"), os.Getenv("TENANT_DB")}

	if DbOption == 1 {
		postgresDb, err := gorm.Open(postgres.Open(os.Getenv("POSTGRES_DB_URL")), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect to PostgreSQL server")
		}
		for _, db := range dbNames {
			CheckAndCreateDatabase(postgresDb, db)
		}
	} else if DbOption == 2 {
		fmt.Print("1 : Admin\n", "2 : Tenant\n", "Choose a DB to migrate : ")
		var dbChoice int
		fmt.Scan(&dbChoice)
		dsn := fmt.Sprintf("%s%s", os.Getenv("POSTGRES_DB_URL"), os.Getenv("ADMIN_DB"))
		adminDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("failed to connect admin database")
		}
		if dbChoice == 1 {
			if err := migrations.MigrateAdminDB(adminDb).Migrate(); err != nil {
				log.Fatal("Failed to migrate admin database")
			}
			fmt.Println("Admin database migrated successfully")
		} else if dbChoice == 2 {
			var tenant []tenant.Tenant
			if tx := adminDb.Table("tenants").Find(&tenant); tx.Error != nil {
				log.Println("Error: ", tx.Error)
			}
			if len(tenant) == 0 {
				log.Fatal("No tenant found.")
			}
			dsn := fmt.Sprintf("%s%s", os.Getenv("POSTGRES_DB_URL"), os.Getenv("TENANT_DB"))
			tenantDb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				log.Fatal("failed to connect tenant database")
			}
			for _, t := range tenant {
				tenantDb.Exec(fmt.Sprintf(`set search_path="%s"`, t.Organization))
				if err := migrations.MigrateTenantDB(tenantDb).Migrate(); err != nil {
					log.Fatal("Failed to migrate tenant database for tenant ID: ", t.ID)
				}
				fmt.Println("Tenant database migrated successfully for tenant ID: ", t.ID)
			}
		} else {
			fmt.Println("Enter the valid option.")
			return
		}
	} else {
		fmt.Println("Enter the valid option.")
		return
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
