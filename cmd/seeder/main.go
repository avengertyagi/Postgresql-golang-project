package main

import (
	"log"

	"github.com/akshit_tyagi/postgresql_project/src/config"
	"github.com/akshit_tyagi/postgresql_project/src/database/seeders"
	"github.com/akshit_tyagi/postgresql_project/src/models"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	if err := config.InitializeDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	if err := models.AutoMigrate(); err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}
	log.Println("Running seeders...")
	seeders.PermissionSeeder()
	seeders.AdminSeeder()
	log.Println("All Seeders completed successfully!")
}
