package repositories

import (
	"context"
	"errors"

	usermodel "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	staffconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/constants"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/dto"
	"gorm.io/gorm"
)

type StaffRepository interface {
	Create(ctx context.Context, staff *usermodel.User) error
	FindByID(ctx context.Context, id uint) (*usermodel.User, error)
	FindByName(ctx context.Context, name string) (*usermodel.User, error)
	FindAll(ctx context.Context, q dto.PaginationQuery) ([]usermodel.User, int64, error)
	Update(ctx context.Context, staff *usermodel.User) error
	Delete(ctx context.Context, id uint) error
	ExistsByName(ctx context.Context, name string, excludeID uint) (bool, error)
}

type staffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) StaffRepository {
	return &staffRepository{db: db}
}

func (r *staffRepository) Create(ctx context.Context, user *usermodel.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *staffRepository) FindAll(ctx context.Context, query dto.PaginationQuery) ([]usermodel.User, int64, error) {
	var total int64
	var users []usermodel.User
	staffQuery := r.db.WithContext(ctx).Model(&usermodel.User{})
	if query.Search != "" {
		staffQuery = staffQuery.Where("LOWER(name) LIKE LOWER(?)", "%"+query.Search+"%")
	}
	if query.Status != "" {
		staffQuery = staffQuery.Where("status = ?", query.Status)
	}
	if err := staffQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (query.Page - 1) * query.Limit
	orderClause := query.SortBy + " " + query.SortDir
	if err := staffQuery.Order(orderClause).Offset(offset).Limit(query.Limit).Preload("Role").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil

}

func (r *staffRepository) Update(ctx context.Context, user *usermodel.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *staffRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&usermodel.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return staffconstants.StaffNotFound
	}
	return nil

}

func (r *staffRepository) FindByID(ctx context.Context, id uint) (*usermodel.User, error) {
	var user usermodel.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, staffconstants.StaffNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (r *staffRepository) FindByName(ctx context.Context, name string) (*usermodel.User, error) {
	var user usermodel.User
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, staffconstants.StaffNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *staffRepository) ExistsByName(ctx context.Context, name string, excludeID uint) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&usermodel.User{}).Where("name = ?", name)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
