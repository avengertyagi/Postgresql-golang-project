package seeders

import (
	"log"
	"time"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	permissionmodel "github.com/akshit_tyagi/postgresql_project/internal/models/permission"
)

func PermissionSeeder() {
	permissions := []permissionmodel.Permission{
		{
			Name:      "role-list",
			GuardName: "web",
			Status:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "role-create",
			GuardName: "web",
			Status:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "role-update",
			GuardName: "web",
			Status:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "role-edit",
			GuardName: "web",
			Status:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:      "role-delete",
			GuardName: "web",
			Status:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	for _, permission := range permissions {
		result := config.DB.Where(permissionmodel.Permission{Name: permission.Name}).FirstOrCreate(&permission)
		if result.Error != nil {
			log.Printf("Failed to seed permission %s: %v", permission.Name, result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("Seeded permission: %s", permission.Name)
		} else {
			log.Printf("Permission already exists, skipped: %s", permission.Name)
		}
	}
}
