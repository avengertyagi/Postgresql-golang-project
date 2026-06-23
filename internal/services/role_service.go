package services

import (
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/akshit_tyagi/postgresql_project/internal/models"
	"github.com/akshit_tyagi/postgresql_project/internal/repositories"
)

func CreateRole(req models.RoleRequest) (*models.Role, error) {
	existing, err := repositories.FindRoleByName(req.Name)
	if err == nil && existing != nil {
		return nil, constants.RoleAlreadyExists
	}
	role := &models.Role{
		Name:   req.Name,
		Status: true,
	}
	if err := repositories.SaveRole(role); err != nil {
		return nil, err
	}
	if len(req.PermissionIDs) > 0 {
		if err := repositories.SyncRolePermissions(role, req.PermissionIDs); err != nil {
			return nil, err
		}
	}
	return role, nil
}
