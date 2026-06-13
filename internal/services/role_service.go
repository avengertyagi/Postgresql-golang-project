package services

import (
	"errors"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/akshit_tyagi/postgresql_project/internal/models"
)

func CreateRole(req models.RoleRequest) (*models.Role, error) {
	existing, err := models.FindRoleByName(req.Name)
	if err == nil && existing != nil {
		return nil, errors.New(constants.RoleAlreadyExists)
	}
	role := &models.Role{
		Name:   req.Name,
		Status: true,
	}
	if err := models.SaveRole(role); err != nil {
		return nil, err
	}
	return role, nil
}
