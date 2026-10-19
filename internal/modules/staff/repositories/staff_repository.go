package repositories

import (
	"context"
	"errors"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	usermodel "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	staffconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/constants"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/dto"
	"gorm.io/gorm"
)

func Create(ctx context.Context, user *usermodel.User) error {
	return config.DB.WithContext(ctx).Create(user).Error
}

func FindAll(ctx context.Context, query dto.PaginationQuery) ([]usermodel.User, int64, error) {
	var total int64
	var users []usermodel.User
	staffQuery := config.DB.WithContext(ctx).Model(&usermodel.User{})
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

func Update(ctx context.Context, user *usermodel.User) error {
	return config.DB.WithContext(ctx).Save(user).Error
}

func Delete(ctx context.Context, id uint) error {
	result := config.DB.WithContext(ctx).Delete(&usermodel.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return staffconstants.StaffNotFound
	}
	return nil
}

func FindByID(ctx context.Context, id uint) (*usermodel.User, error) {
	var user usermodel.User
	err := config.DB.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, staffconstants.StaffNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FindByName(ctx context.Context, name string) (*usermodel.User, error) {
	var user usermodel.User
	err := config.DB.WithContext(ctx).Where("name = ?", name).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, staffconstants.StaffNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func ExistsByMobile(ctx context.Context, mobile int, excludeID uint) (bool, error) {
	var count int64
	q := config.DB.WithContext(ctx).Model(&usermodel.User{}).Where("mobile = ?", mobile)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
