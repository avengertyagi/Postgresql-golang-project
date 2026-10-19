package repositories

import (
	"context"
	"errors"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	themeconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/constants"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/dto"
	thememodel "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/models"
	"gorm.io/gorm"
)

func Create(ctx context.Context, theme *thememodel.Theme) error {
	return config.DB.WithContext(ctx).Create(theme).Error
}

func FindAll(ctx context.Context, query dto.PaginationQuery) ([]thememodel.Theme, int64, error) {
	var total int64
	var themes []thememodel.Theme
	theme := config.DB.WithContext(ctx).Model(&thememodel.Theme{})
	if query.Search != "" {
		theme = theme.Where("LOWER(name) LIKE LOWER(?)", "%"+query.Search+"%")
	}
	if query.Status != "" {
		theme = theme.Where("status = ?", query.Status)
	}
	if err := theme.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (query.Page - 1) * query.Limit
	orderClause := query.SortBy + " " + query.SortDir
	if err := theme.Order(orderClause).Offset(offset).Limit(query.Limit).Find(&themes).Error; err != nil {
		return nil, 0, err
	}
	return themes, total, nil
}

func Update(ctx context.Context, theme *thememodel.Theme) error {
	return config.DB.WithContext(ctx).Save(theme).Error
}

func Delete(ctx context.Context, id uint) error {
	result := config.DB.WithContext(ctx).Delete(&thememodel.Theme{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return themeconstants.ThemeNotFound
	}
	return nil
}

func FindByID(ctx context.Context, id uint) (*thememodel.Theme, error) {
	var theme thememodel.Theme
	err := config.DB.WithContext(ctx).First(&theme, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, themeconstants.ThemeNotFound
	}
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

func FindByName(ctx context.Context, name string) (*thememodel.Theme, error) {
	var theme thememodel.Theme
	err := config.DB.WithContext(ctx).Where("name = ?", name).First(&theme).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, themeconstants.ThemeNotFound
	}
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

func ExistsByName(ctx context.Context, name string, excludeID uint) (bool, error) {
	var count int64
	q := config.DB.WithContext(ctx).Model(&thememodel.Theme{}).Where("name = ?", name)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
