package seeders

import (
	"fmt"
	"log"
	"os"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	usermodel "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	permissionmodel "github.com/akshit_tyagi/postgresql_project/internal/modules/permission/models"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"

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
	fmt.Println("constants.SuperAdminRole", constants.SuperAdminRole)
	var permissions []permissionmodel.Permission
	if err := config.DB.Find(&permissions).Error; err != nil {
		log.Fatalf("Failed to fetch permissions for seeding: %v", err)
	}
	var adminRole rolemodel.Role
	if err := config.DB.Where(rolemodel.Role{Name: "Admin"}).FirstOrCreate(&adminRole).Error; err != nil {
		log.Fatalf("Failed to find or create Admin role: %v", err)
	}
	if err := config.DB.Model(&adminRole).Association("Permissions").Replace(permissions); err != nil {
		log.Fatalf("Failed to sync permissions to Admin role: %v", err)
	}
	log.Printf("Synced %d permissions to Admin role.", len(permissions))

	admins := []usermodel.User{
		{
			Name:     "Super Admin",
			Email:    "superadmin@gmail.com",
			Password: string(hashedPassword),
			Status:   true,
			UserType: constants.SuperAdminRole,
			RoleID:   adminRole.ID,
		},
	}
	for _, admin := range admins {
		admin.RoleID = adminRole.ID
		admin.UserType = constants.SuperAdminRole

		result := config.DB.Where("email = ?", admin.Email).FirstOrCreate(&admin)
		if result.Error != nil {
			log.Printf("Failed to seed admin %s: %v", admin.Email, result.Error)
			continue
		}

		if err := config.DB.Model(&admin).Updates(map[string]interface{}{
			"user_type": constants.SuperAdminRole,
			"role_id":   adminRole.ID,
		}).Error; err != nil {
			log.Printf("Failed to update admin fields for %s: %v", admin.Email, err)
		}

		if result.RowsAffected > 0 {
			log.Printf("Seeded admin: %s", admin.Email)
		} else {
			log.Printf("Admin already exists, updated role and user_type: %s", admin.Email)
		}
	}
}
