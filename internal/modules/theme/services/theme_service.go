package services

import (
	"context"
	"strings"

	themeconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/constants"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/dto"
	thememodel "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/models"
	themerepo "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/repositories"
)

type ThemeService interface {
	Create(ctx context.Context, req dto.ThemeRequest) (*thememodel.Theme, error)
	GetByID(ctx context.Context, id uint) (*thememodel.Theme, error)
	List(ctx context.Context, q dto.PaginationQuery) ([]thememodel.Theme, int64, error)
	Update(ctx context.Context, id uint, req dto.ThemeRequest) (*thememodel.Theme, error)
	UpdateStatus(ctx context.Context, id uint) (*thememodel.Theme, error)
	Delete(ctx context.Context, id uint) (*thememodel.Theme, error)
}

type themeService struct {
	repo themerepo.ThemeRepository
}

func NewThemeService(repo themerepo.ThemeRepository) ThemeService {
	return &themeService{repo: repo}
}

func (s *themeService) List(ctx context.Context, q dto.PaginationQuery) ([]thememodel.Theme, int64, error) {
	return s.repo.FindAll(ctx, q)
}

func (s *themeService) Create(ctx context.Context, req dto.ThemeRequest) (*thememodel.Theme, error) {
	name := strings.TrimSpace(req.Name)
	exists, err := s.repo.ExistsByName(ctx, name, 0)
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
	if err := s.repo.Create(ctx, theme); err != nil {
		return nil, err
	}
	return theme, nil
}

func (s *themeService) GetByID(ctx context.Context, id uint) (*thememodel.Theme, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *themeService) Update(ctx context.Context, id uint, req dto.ThemeRequest) (*thememodel.Theme, error) {
	theme, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if !strings.EqualFold(name, theme.Name) {
		exists, err := s.repo.ExistsByName(ctx, name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, themeconstants.ThemeAlreadyExists
		}
		theme.Name = name
	}
	if err := s.repo.Update(ctx, theme); err != nil {
		return nil, err
	}

	return theme, nil
}

func (s *themeService) UpdateStatus(ctx context.Context, id uint) (*thememodel.Theme, error) {
	theme, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	theme.Status = !theme.Status
	if err := s.repo.Update(ctx, theme); err != nil {
		return nil, err
	}
	return theme, nil
}

func (s *themeService) Delete(ctx context.Context, id uint) (*thememodel.Theme, error) {
	theme, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return nil, err
	}
	return theme, nil
}
