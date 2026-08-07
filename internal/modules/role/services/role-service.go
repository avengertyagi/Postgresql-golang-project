package role

import (
	"errors"
	"math"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	permissionmodel "github.com/akshit_tyagi/postgresql_project/internal/modules/permission/models"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
)

func GetAll(req rolemodel.RoleListRequest) (*rolemodel.RoleListResponse, error) {
	if req.CurrentPage <= 0 {
		req.CurrentPage = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 10
	}
	var total int64
	var roleList []rolemodel.Role
	query := config.DB.Model(&rolemodel.Role{})
	if req.Search != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+req.Search+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	if req.SortBy != "name" && req.SortBy != "created_at" && req.SortBy != "updated_at" {
		req.SortBy = "id"
	}
	if req.SortOrder != "desc" && req.SortOrder != "DESC" {
		req.SortOrder = "ASC"
	}
	offset := (req.CurrentPage - 1) * req.PerPage
	if err := query.
		Order(req.SortBy + " " + req.SortOrder).
		Offset(offset).
		Limit(req.PerPage).
		Find(&roleList).Error; err != nil {
		return nil, err
	}
	lastPage := int(math.Ceil(float64(total) / float64(req.PerPage)))
	return &rolemodel.RoleListResponse{
		Data:        roleList,
		CurrentPage: req.CurrentPage,
		PerPage:     req.PerPage,
		Total:       total,
		LastPage:    lastPage,
	}, nil
}

func Create(req rolemodel.RoleRequest) (*rolemodel.Role, error) {
	existing, err := FindRoleByName(req.Name)
	if err == nil && existing != nil {
		return nil, constants.RoleAlreadyExists
	}
	role := &rolemodel.Role{
		Name:   req.Name,
		Status: true,
	}
	if err := createRole(role); err != nil {
		return nil, err
	}
	if len(req.PermissionIDs) > 0 {
		if err := SyncRolePermissions(role, req.PermissionIDs); err != nil {
			return nil, err
		}
	}
	return role, nil
}

func GetByID(ID string) (*rolemodel.Role, error) {
	role, err := FindByID(ID)
	if err != nil {
		if errors.Is(err, constants.RoleNotFound) {
			return nil, constants.RoleNotFound
		}
		return nil, err
	}
	return role, nil
}

func Update(id string, req rolemodel.RoleRequest) (*rolemodel.Role, error) {
	role, err := FindByID(id)
	if err != nil {
		return nil, constants.RoleNotFound
	}
	existing, err := FindRoleByName(req.Name, role.ID)
	if err == nil && existing != nil {
		return nil, constants.RoleAlreadyExists
	}
	role.Name = req.Name
	if err := updateRole(role); err != nil {
		return nil, err
	}
	if err := SyncRolePermissions(role, req.PermissionIDs); err != nil {
		return nil, err
	}
	return role, nil
}

func Delete(id string) (*rolemodel.Role, error) {
	role, err := FindByID(id)
	if err != nil {
		return nil, constants.RoleNotFound
	}
	if err := deleteRole(role); err != nil {
		return nil, err
	}
	return role, nil
}

func createRole(role *rolemodel.Role) error {
	return config.DB.Create(role).Error
}

func updateRole(role *rolemodel.Role) error {
	return config.DB.Save(role).Error
}

func deleteRole(role *rolemodel.Role) error {
	return config.DB.Delete(role).Error
}

func FindByID(id string) (*rolemodel.Role, error) {
	var role rolemodel.Role
	err := config.DB.Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func FindRoleByName(name string, excludeID ...uint) (*rolemodel.Role, error) {
	var role rolemodel.Role
	query := config.DB.Where("LOWER(name) = LOWER(?)", name)
	if len(excludeID) > 0 {
		query = query.Where("id <> ?", excludeID[0])
	}
	err := query.First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
func SyncRolePermissions(role *rolemodel.Role, permissionIDs []uint) error {
	var permissions []permissionmodel.Permission
	if err := config.DB.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}
	return config.DB.Model(role).Association("Permissions").Replace(permissions)
}
