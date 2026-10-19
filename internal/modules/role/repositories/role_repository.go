package repositories

import (
	"context"
	"errors"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	permissionmodel "github.com/akshit_tyagi/postgresql_project/internal/modules/permission/models"
	roleconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/role/constants"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/role/dto"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
	"gorm.io/gorm"
)

func Create(ctx context.Context, role *rolemodel.Role) error {
	return config.DB.WithContext(ctx).Create(role).Error
}

func FindAll(ctx context.Context, query dto.PaginationQuery) ([]rolemodel.Role, int64, error) {
	var total int64
	var roles []rolemodel.Role
	role := config.DB.WithContext(ctx).Model(&rolemodel.Role{})
	if query.Search != "" {
		role = role.Where("LOWER(name) LIKE LOWER(?)", "%"+query.Search+"%")
	}
	if query.Status != "" {
		role = role.Where("status = ?", query.Status)
	}
	if err := role.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (query.Page - 1) * query.Limit
	orderClause := query.SortBy + " " + query.SortDir
	if err := role.Order(orderClause).Offset(offset).Limit(query.Limit).Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil
}

func Update(ctx context.Context, role *rolemodel.Role) error {
	return config.DB.WithContext(ctx).Save(role).Error
}

func Delete(ctx context.Context, id uint) error {
	result := config.DB.WithContext(ctx).Delete(&rolemodel.Role{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return roleconstants.RoleNotFound
	}
	return nil
}

func FindByID(ctx context.Context, id uint) (*rolemodel.Role, error) {
	var role rolemodel.Role
	err := config.DB.WithContext(ctx).Preload("Permissions").First(&role, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, roleconstants.RoleNotFound
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func FindByName(ctx context.Context, name string) (*rolemodel.Role, error) {
	var role rolemodel.Role
	err := config.DB.WithContext(ctx).Where("name = ?", name).Preload("Permissions").First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, roleconstants.RoleNotFound
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func ExistsByName(ctx context.Context, name string, excludeID uint) (bool, error) {
	var count int64
	q := config.DB.WithContext(ctx).Model(&rolemodel.Role{}).Where("name = ?", name)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func SyncPermissions(ctx context.Context, role *rolemodel.Role, permissionIDs []uint) error {
	var permissions []permissionmodel.Permission
	if err := config.DB.WithContext(ctx).Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}
	return config.DB.WithContext(ctx).Model(role).Association("Permissions").Replace(permissions)
}
