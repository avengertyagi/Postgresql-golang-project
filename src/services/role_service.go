package service

import (
	"errors"

	"github.com/akshit_tyagi/postgresql_project/src/constants"
	authmodel "github.com/akshit_tyagi/postgresql_project/src/models"
)

func CreateRole(req authmodel.RoleRequest) (*authmodel.Role, error) {
	existing, err := authmodel.FindRoleByName(req.Name)
	if err == nil && existing != nil {
		return nil, errors.New(constants.ROLE_ALREADY_EXISTS)
	}
	role := &authmodel.Role{
		Name:   req.Name,
		Status: true,
	}
	if err := authmodel.SaveRole(role); err != nil {
		return nil, err
	}
	return role, nil
}
