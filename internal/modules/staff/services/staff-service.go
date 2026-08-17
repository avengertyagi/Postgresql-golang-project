package services

import (
	"context"
	"strings"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	usermodel "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/dto"
	staffrepo "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/repositories"
)

type StaffService interface {
	Create(ctx context.Context, req dto.StaffRequest) (*usermodel.User, error)
	GetByID(ctx context.Context, id uint) (*usermodel.User, error)
	List(ctx context.Context, q dto.PaginationQuery) ([]usermodel.User, int64, error)
	Update(ctx context.Context, id uint, req dto.StaffRequest) (*usermodel.User, error)
	Delete(ctx context.Context, id uint) (*usermodel.User, error)
}

type staffService struct {
	repo staffrepo.StaffRepository
}

func NewStaffService(repo staffrepo.StaffRepository) StaffService {
	return &staffService{repo: repo}
}

func (s *staffService) List(ctx context.Context, q dto.PaginationQuery) ([]usermodel.User, int64, error) {
	return s.repo.FindAll(ctx, q)
}

func (s *staffService) Create(ctx context.Context, req dto.StaffRequest) (*usermodel.User, error) {
	name := strings.TrimSpace(req.Name)
	exists, err := s.repo.ExistsByName(ctx, name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, constants.StaffAlreadyExists
	}
	staff := &usermodel.User{
		Name:     name,
		Email:    req.Email,
		RoleID:   req.RoleID,
		Mobile:   req.Mobile,
		Password: req.Password,
		UserType: 1,
	}
	if err := s.repo.Create(ctx, staff); err != nil {
		return nil, err
	}
	return staff, nil
}

func (s *staffService) GetByID(ctx context.Context, id uint) (*usermodel.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *staffService) Update(ctx context.Context, id uint, req dto.StaffRequest) (*usermodel.User, error) {
	staff, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if !strings.EqualFold(name, staff.Name) {
		exists, err := s.repo.ExistsByName(ctx, name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, constants.StaffAlreadyExists
		}
		staff.Name = name
	}
	if err := s.repo.Update(ctx, staff); err != nil {
		return nil, err
	}

	return staff, nil
}

func (s *staffService) Delete(ctx context.Context, id uint) (*usermodel.User, error) {
	staff, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return nil, err
	}
	return staff, nil
}
