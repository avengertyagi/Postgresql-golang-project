package services

import (
	"errors"
	"math"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/models/role"
	rolerepo "github.com/akshit_tyagi/postgresql_project/internal/repositories/role"
)

func GetAll(req rolemodel.RoleListRequest) (*rolemodel.RoleListResponse, error) {
	if req.CurrentPage <= 0 {
		req.CurrentPage = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 10
	}
	roles, total, err := rolerepo.FindAll(req)
	if err != nil {
		return nil, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(req.PerPage)))
	return &rolemodel.RoleListResponse{
		Data:        roles,
		CurrentPage: req.CurrentPage,
		PerPage:     req.PerPage,
		Total:       total,
		LastPage:    lastPage,
	}, nil
}

func Create(req rolemodel.RoleRequest) (*rolemodel.Role, error) {
	existing, err := rolerepo.FindRoleByName(req.Name)
	if err == nil && existing != nil {
		return nil, constants.RoleAlreadyExists
	}
	role := &rolemodel.Role{
		Name:   req.Name,
		Status: true,
	}
	if err := rolerepo.Create(role); err != nil {
		return nil, err
	}
	if len(req.PermissionIDs) > 0 {
		if err := rolerepo.SyncRolePermissions(role, req.PermissionIDs); err != nil {
			return nil, err
		}
	}
	return role, nil
}

func GetByID(ID string) (*rolemodel.Role, error) {
	role, err := rolerepo.FindByID(ID)
	if err != nil {
		if errors.Is(err, constants.RoleNotFound) {
			return nil, constants.RoleNotFound
		}
		return nil, err
	}
	return role, nil
}

func Update(id string, req rolemodel.RoleRequest) (*rolemodel.Role, error) {
	role, err := rolerepo.FindByID(id)
	if err != nil {
		return nil, constants.RoleNotFound
	}
	existing, err := rolerepo.FindRoleByName(req.Name, role.ID)
	if err == nil && existing != nil {
		return nil, constants.RoleAlreadyExists
	}
	role.Name = req.Name
	if err := rolerepo.Update(role); err != nil {
		return nil, err
	}
	if err := rolerepo.SyncRolePermissions(role, req.PermissionIDs); err != nil {
		return nil, err
	}
	return role, nil
}

func Delete(id string) (*rolemodel.Role, error) {
	role, err := rolerepo.FindByID(id)
	if err != nil {
		return nil, constants.RoleNotFound
	}
	if err := rolerepo.Delete(role); err != nil {
		return nil, err
	}
	return role, nil
}
