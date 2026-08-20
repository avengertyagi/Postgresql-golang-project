package seeders

import (
	"log"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	permissionmodel "github.com/akshit_tyagi/postgresql_project/internal/modules/permission/models"
)

func PermissionSeeder() {
	permissions := []permissionmodel.Permission{
		{
			Name:      "role-list",
			GuardName: "web",
		},
		{
			Name:      "role-create",
			GuardName: "web",
		},
		{
			Name:      "role-update",
			GuardName: "web",
		},
		{
			Name:      "role-edit",
			GuardName: "web",
		},
		{
			Name:      "role-delete",
			GuardName: "web",
		},
		{
			Name:      "staff-list",
			GuardName: "web",
		},
		{
			Name:      "staff-create",
			GuardName: "web",
		},
		{
			Name:      "staff-update",
			GuardName: "web",
		},
		{
			Name:      "staff-edit",
			GuardName: "web",
		},
		{
			Name:      "staff-delete",
			GuardName: "web",
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
