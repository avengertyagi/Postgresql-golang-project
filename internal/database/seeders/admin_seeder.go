package seeders

import (
	"log"
	"os"
	"time"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/akshit_tyagi/postgresql_project/internal/models"
	"github.com/akshit_tyagi/postgresql_project/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

func AdminSeeder() {
	seedPassword := os.Getenv("ADMIN_SEEDER_PASSWORD")
	if seedPassword == "" {
		log.Fatal("ADMIN_SEEDER_PASSWORD env var is required for seeding. Set it in .env file.")
	}
	if len(seedPassword) < 8 {
		log.Fatal("ADMIN_SEEDER_PASSWORD must be at least 8 characters long.")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(seedPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Fatal(err)
	}

	admins := []models.User{
		{
			Name:      "Super Admin",
			Email:     "superadmin@gmail.com",
			Password:  string(hashedPassword),
			Status:    true,
			UserType:  constants.SuperAdminRole,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	var permissions []models.Permission
	if err := config.DB.Find(&permissions).Error; err != nil {
		log.Fatalf("Failed to fetch permissions for seeding: %v", err)
	}
	var adminRole models.Role
	if err := config.DB.Where(models.Role{Name: "Admin"}).FirstOrCreate(&adminRole).Error; err != nil {
		log.Fatalf("Failed to find or create Admin role: %v", err)
	}
	if err := adminRole.SyncPermissions(permissions); err != nil {
		log.Fatalf("Failed to sync permissions to Admin role: %v", err)
	}
	log.Printf("Synced %d permissions to Admin role.", len(permissions))
	for _, admin := range admins {
		result := config.DB.Where(models.User{Email: admin.Email}).FirstOrCreate(&admin)
		if result.Error != nil {
			log.Printf("Failed to seed admin %s: %v", admin.Email, result.Error)
			continue
		}
		if result.RowsAffected > 0 {
			log.Printf("Seeded admin: %s", admin.Email)
		} else {
			log.Printf("Admin already exists, skipped: %s", admin.Email)
		}
		if err := repositories.AssignRole(&admin, &adminRole); err != nil {
			log.Printf("Failed to assign Admin role to %s: %v", admin.Email, err)
		} else {
			log.Printf("Assigned Admin role to %s", admin.Email)
		}
	}
}
