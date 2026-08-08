package role

import (
	"errors"
	"math"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
	rolerepo "github.com/akshit_tyagi/postgresql_project/internal/modules/role/repositories"
)

func GetAllWithPagination(req rolemodel.RoleListRequest) ([]rolemodel.Role, int64, int, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.SortBy != "name" && req.SortBy != "created_at" && req.SortBy != "updated_at" {
		req.SortBy = "id"
	}
	if req.SortOrder != "desc" && req.SortOrder != "DESC" {
		req.SortOrder = "ASC"
	}
	roleList, total, err := rolerepo.ListRole(req)
	if err != nil {
		return nil, 0, 0, err
	}
	lastPage := int(math.Ceil(float64(total) / float64(req.Limit)))
	return roleList, total, lastPage, nil
}

func Create(req rolemodel.RoleRequest) (*rolemodel.Role, error) {
	existing, err := rolerepo.FindRoleByName(req.Name)
	if err == nil && existing != nil {
		return nil, constants.RoleAlreadyExists
	}
	role, err := rolerepo.CreateRole(&req)
	if err != nil {
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
	if err := rolerepo.UpdateRole(role); err != nil {
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
	if err := rolerepo.DeleteRole(role); err != nil {
		return nil, err
	}
	return role, nil
}
