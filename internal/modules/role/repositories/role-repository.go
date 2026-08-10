package repositories

import (
	"context"
	"errors"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	permissionmodel "github.com/akshit_tyagi/postgresql_project/internal/modules/permission/models"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/role/dto"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(ctx context.Context, role *rolemodel.Role) error
	FindByID(ctx context.Context, id uint) (*rolemodel.Role, error)
	FindByName(ctx context.Context, name string) (*rolemodel.Role, error)
	FindAll(ctx context.Context, q dto.PaginationQuery) ([]rolemodel.Role, int64, error)
	Update(ctx context.Context, role *rolemodel.Role) error
	Delete(ctx context.Context, id uint) error
	ExistsByName(ctx context.Context, name string, excludeID uint) (bool, error)
	SyncPermissions(ctx context.Context, role *rolemodel.Role, permissionIDs []uint) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *rolemodel.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *roleRepository) FindAll(ctx context.Context, query dto.PaginationQuery) ([]rolemodel.Role, int64, error) {
	var total int64
	var roles []rolemodel.Role
	role := r.db.WithContext(ctx).Model(&rolemodel.Role{})
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
	if err := role.Order(orderClause).Offset(offset).Limit(query.Limit).Find(&roles).Error; err != nil {
		return nil, 0, err
	}
	return roles, total, nil

}

func (r *roleRepository) Update(ctx context.Context, role *rolemodel.Role) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *roleRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&rolemodel.Role{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return constants.RoleNotFound
	}
	return nil

}

func (r *roleRepository) FindByID(ctx context.Context, id uint) (*rolemodel.Role, error) {
	var role rolemodel.Role
	err := r.db.WithContext(ctx).First(&role, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, constants.RoleNotFound
	}
	if err != nil {
		return nil, err
	}
	return &role, nil

}

func (r *roleRepository) FindByName(ctx context.Context, name string) (*rolemodel.Role, error) {
	var role rolemodel.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, constants.RoleNotFound
	}
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) ExistsByName(ctx context.Context, name string, excludeID uint) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&rolemodel.Role{}).Where("name = ?", name)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *roleRepository) SyncPermissions(ctx context.Context, role *rolemodel.Role, permissionIDs []uint) error {
	var permissions []permissionmodel.Permission
	if err := r.db.WithContext(ctx).Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Model(role).Association("Permissions").Replace(permissions)
}
