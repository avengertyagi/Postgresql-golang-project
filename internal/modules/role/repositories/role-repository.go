package repositories

import (
	"github.com/akshit_tyagi/postgresql_project/internal/config"
	permissionmodel "github.com/akshit_tyagi/postgresql_project/internal/modules/permission/models"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
)

func ListRole(req rolemodel.RoleListRequest) ([]rolemodel.Role, int64, error) {
	var total int64
	var roles []rolemodel.Role
	query := config.DB.Model(&rolemodel.Role{})
	if req.Search != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+req.Search+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (req.Page - 1) * req.Limit
	if err := query.
		Preload("Permissions").
		Order(req.SortBy + " " + req.SortOrder).
		Offset(offset).
		Limit(req.Limit).
		Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}
func CreateRole(req *rolemodel.RoleRequest) (*rolemodel.Role, error) {
	role := &rolemodel.Role{
		Name:   req.Name,
		Status: true,
	}
	if err := config.DB.Model(&rolemodel.Role{}).Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func UpdateRole(role *rolemodel.Role) error {
	return config.DB.Save(role).Error
}

func DeleteRole(role *rolemodel.Role) error {
	return config.DB.Delete(role).Error
}

func FindByID(id string) (*rolemodel.Role, error) {
	var role rolemodel.Role
	err := config.DB.Model(&rolemodel.Role{}).Preload("Permissions").Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func FindRoleByName(name string, excludeID ...uint) (*rolemodel.Role, error) {
	var role rolemodel.Role
	query := config.DB.Model(&rolemodel.Role{}).Where("LOWER(name) = LOWER(?)", name)
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
