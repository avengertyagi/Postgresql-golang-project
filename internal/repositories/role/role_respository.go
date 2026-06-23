package repositories

import (
	"github.com/akshit_tyagi/postgresql_project/internal/config"
	permissionmodel "github.com/akshit_tyagi/postgresql_project/internal/models/permission"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/models/role"
)

func Create(role *rolemodel.Role) error {
	return config.DB.Create(role).Error
}

func Update(role *rolemodel.Role) error {
	return config.DB.Save(role).Error
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

func FindAll(req rolemodel.RoleListRequest) ([]rolemodel.Role, int64, error) {
	var roles []rolemodel.Role
	var total int64
	query := config.DB.Model(&rolemodel.Role{})
	if req.Search != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+req.Search+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if req.SortBy != "name" && req.SortBy != "created_at" && req.SortBy != "updated_at" {
		req.SortBy = "id"
	}
	if req.SortOrder != "desc" && req.SortOrder != "DESC" {
		req.SortOrder = "ASC"
	}
	offset := (req.CurrentPage - 1) * req.PerPage
	err := query.
		Order(req.SortBy + " " + req.SortOrder).
		Offset(offset).
		Limit(req.PerPage).
		Find(&roles).Error

	if err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func Delete(role *rolemodel.Role) error {
	return config.DB.Delete(role).Error
}
