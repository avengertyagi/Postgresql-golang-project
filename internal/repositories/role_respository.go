package repositories

import (
	"github.com/akshit_tyagi/postgresql_project/internal/config"
	"github.com/akshit_tyagi/postgresql_project/internal/models"
)

func SaveRole(role *models.Role) error {
	return config.DB.Create(role).Error
}

func FindRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := config.DB.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func SyncRolePermissions(role *models.Role, permissionIDs []uint) error {
	var permissions []models.Permission
	if err := config.DB.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}
	return config.DB.Model(role).Association("Permissions").Replace(permissions)
}
