package services

import (
	"context"
	"strings"

	roleconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/role/constants"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/role/dto"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
	rolerepo "github.com/akshit_tyagi/postgresql_project/internal/modules/role/repositories"
)

func List(ctx context.Context, q dto.PaginationQuery) ([]rolemodel.Role, int64, error) {
	return rolerepo.FindAll(ctx, q)
}

func Create(ctx context.Context, req dto.RoleRequest) (*rolemodel.Role, error) {
	name := strings.TrimSpace(req.Name)
	exists, err := rolerepo.ExistsByName(ctx, name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, roleconstants.RoleAlreadyExists
	}
	role := &rolemodel.Role{
		Name: name,
	}
	if err := rolerepo.Create(ctx, role); err != nil {
		return nil, err
	}
	if len(req.PermissionIDs) > 0 {
		if err := rolerepo.SyncPermissions(ctx, role, req.PermissionIDs); err != nil {
			return nil, err
		}
	}
	return role, nil
}

func GetByID(ctx context.Context, id uint) (*rolemodel.Role, error) {
	return rolerepo.FindByID(ctx, id)
}

func Update(ctx context.Context, id uint, req dto.RoleRequest) (*rolemodel.Role, error) {
	role, err := rolerepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if !strings.EqualFold(name, role.Name) {
		exists, err := rolerepo.ExistsByName(ctx, name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, roleconstants.RoleAlreadyExists
		}
		role.Name = name
	}
	if err := rolerepo.Update(ctx, role); err != nil {
		return nil, err
	}
	if err := rolerepo.SyncPermissions(ctx, role, req.PermissionIDs); err != nil {
		return nil, err
	}

	return role, nil
}

func UpdateStatus(ctx context.Context, id uint) (*rolemodel.Role, error) {
	role, err := rolerepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	role.Status = !role.Status
	if err := rolerepo.Update(ctx, role); err != nil {
		return nil, err
	}
	return role, nil
}

func Delete(ctx context.Context, id uint) (*rolemodel.Role, error) {
	role, err := rolerepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := rolerepo.Delete(ctx, id); err != nil {
		return nil, err
	}
	return role, nil
}
