package repositories

import (
	"context"
	"errors"

	themeconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/constants"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/dto"
	thememodel "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/models"
	"gorm.io/gorm"
)

type ThemeRepository interface {
	Create(ctx context.Context, theme *thememodel.Theme) error
	FindByID(ctx context.Context, id uint) (*thememodel.Theme, error)
	FindByName(ctx context.Context, name string) (*thememodel.Theme, error)
	FindAll(ctx context.Context, q dto.PaginationQuery) ([]thememodel.Theme, int64, error)
	Update(ctx context.Context, theme *thememodel.Theme) error
	Delete(ctx context.Context, id uint) error
	ExistsByName(ctx context.Context, name string, excludeID uint) (bool, error)
}

type themeRepository struct {
	db *gorm.DB
}

func NewThemeRepository(db *gorm.DB) ThemeRepository {
	return &themeRepository{db: db}
}

func (r *themeRepository) Create(ctx context.Context, theme *thememodel.Theme) error {
	return r.db.WithContext(ctx).Create(theme).Error
}

func (r *themeRepository) FindAll(ctx context.Context, query dto.PaginationQuery) ([]thememodel.Theme, int64, error) {
	var total int64
	var themes []thememodel.Theme
	theme := r.db.WithContext(ctx).Model(&thememodel.Theme{})
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

func (r *themeRepository) Update(ctx context.Context, theme *thememodel.Theme) error {
	return r.db.WithContext(ctx).Save(theme).Error
}

func (r *themeRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&thememodel.Theme{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return themeconstants.ThemeNotFound
	}
	return nil

}

func (r *themeRepository) FindByID(ctx context.Context, id uint) (*thememodel.Theme, error) {
	var theme thememodel.Theme
	err := r.db.WithContext(ctx).First(&theme, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, themeconstants.ThemeNotFound
	}
	if err != nil {
		return nil, err
	}
	return &theme, nil

}

func (r *themeRepository) FindByName(ctx context.Context, name string) (*thememodel.Theme, error) {
	var theme thememodel.Theme
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&theme).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, themeconstants.ThemeNotFound
	}
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

func (r *themeRepository) ExistsByName(ctx context.Context, name string, excludeID uint) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&thememodel.Theme{}).Where("name = ?", name)
	if excludeID > 0 {
		q = q.Where("id <> ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
