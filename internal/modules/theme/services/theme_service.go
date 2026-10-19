package services

import (
	"context"
	"strings"

	themeconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/constants"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/dto"
	thememodel "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/models"
	themerepo "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/repositories"
)

func List(ctx context.Context, q dto.PaginationQuery) ([]thememodel.Theme, int64, error) {
	return themerepo.FindAll(ctx, q)
}

func Create(ctx context.Context, req dto.ThemeRequest) (*thememodel.Theme, error) {
	name := strings.TrimSpace(req.Name)
	exists, err := themerepo.ExistsByName(ctx, name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, themeconstants.ThemeAlreadyExists
	}
	theme := &thememodel.Theme{
		Name:        name,
		Description: req.Description,
		Image:       req.Image,
	}
	if err := themerepo.Create(ctx, theme); err != nil {
		return nil, err
	}
	return theme, nil
}

func GetByID(ctx context.Context, id uint) (*thememodel.Theme, error) {
	return themerepo.FindByID(ctx, id)
}

func Update(ctx context.Context, id uint, req dto.ThemeRequest) (*thememodel.Theme, error) {
	theme, err := themerepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if !strings.EqualFold(name, theme.Name) {
		exists, err := themerepo.ExistsByName(ctx, name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, themeconstants.ThemeAlreadyExists
		}
		theme.Name = name
	}
	if err := themerepo.Update(ctx, theme); err != nil {
		return nil, err
	}

	return theme, nil
}

func UpdateStatus(ctx context.Context, id uint) (*thememodel.Theme, error) {
	theme, err := themerepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	theme.Status = !theme.Status
	if err := themerepo.Update(ctx, theme); err != nil {
		return nil, err
	}
	return theme, nil
}

func Delete(ctx context.Context, id uint) (*thememodel.Theme, error) {
	theme, err := themerepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := themerepo.Delete(ctx, id); err != nil {
		return nil, err
	}
	return theme, nil
}
