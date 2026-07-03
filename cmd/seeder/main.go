package main

import (
	"log"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	"github.com/akshit_tyagi/postgresql_project/internal/database/seeders"
	permissionmodel "github.com/akshit_tyagi/postgresql_project/internal/models/permission"
	personalaccesstokenmodel "github.com/akshit_tyagi/postgresql_project/internal/models/personalaccesstoken"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/models/role"
	tenantmodel "github.com/akshit_tyagi/postgresql_project/internal/models/tenant"
	usermodel "github.com/akshit_tyagi/postgresql_project/internal/models/user"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := config.InitializeDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	if err := config.Migrate(
		&permissionmodel.Permission{},
		&rolemodel.Role{},
		&usermodel.User{},
		&tenantmodel.Tenant{},
		&personalaccesstokenmodel.PersonalAccessToken{},
	); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Running seeders...")
	seeders.PermissionSeeder()
	seeders.AdminSeeder()
	log.Println("All Seeders completed successfully!")
}
