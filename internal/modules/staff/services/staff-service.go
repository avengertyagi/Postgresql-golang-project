package services

import (
	"context"
	"strings"

	usermodel "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	rolerepo "github.com/akshit_tyagi/postgresql_project/internal/modules/role/repositories"
	staffconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/staff/constants"
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
	staffRepo staffrepo.StaffRepository
	roleRepo  rolerepo.RoleRepository
}

func NewStaffService(staffRepo staffrepo.StaffRepository, roleRepo rolerepo.RoleRepository) StaffService {
	return &staffService{staffRepo: staffRepo, roleRepo: roleRepo}
}

func (s *staffService) List(ctx context.Context, q dto.PaginationQuery) ([]usermodel.User, int64, error) {
	return s.staffRepo.FindAll(ctx, q)
}

func (s *staffService) Create(ctx context.Context, req dto.StaffRequest) (*usermodel.User, error) {
	name := strings.TrimSpace(req.Name)
	exists, err := s.staffRepo.ExistsByMobile(ctx, req.Mobile, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, staffconstants.StaffAlreadyExists
	}
	role, err := s.roleRepo.FindByID(ctx, req.RoleID)
	if role == nil {
		return nil, staffconstants.RoleIdNotFound
	}

	staff := &usermodel.User{
		Name:     name,
		Email:    req.Email,
		RoleID:   req.RoleID,
		Mobile:   req.Mobile,
		UserType: 1,
	}
	if err := staff.HashPassword(req.Password); err != nil {
		return nil, err
	}
	if err := s.staffRepo.Create(ctx, staff); err != nil {
		return nil, err
	}
	return staff, nil
}

func (s *staffService) GetByID(ctx context.Context, id uint) (*usermodel.User, error) {
	return s.staffRepo.FindByID(ctx, id)
}

func (s *staffService) Update(ctx context.Context, id uint, req dto.StaffRequest) (*usermodel.User, error) {
	staff, err := s.staffRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	exists, err := s.staffRepo.ExistsByMobile(ctx, req.Mobile, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, staffconstants.StaffAlreadyExists
	}
	staff.Name = name
	role, err := s.roleRepo.FindByID(ctx, req.RoleID)
	if role == nil {
		return nil, staffconstants.RoleIdNotFound
	}
	staff.Email = req.Email
	staff.Mobile = req.Mobile
	staff.RoleID = req.RoleID
	if req.Password != "" {
		if err := staff.HashPassword(req.Password); err != nil {
			return nil, err
		}
	}
	if err := s.staffRepo.Update(ctx, staff); err != nil {
		return nil, err
	}

	return staff, nil
}

func (s *staffService) Delete(ctx context.Context, id uint) (*usermodel.User, error) {
	staff, err := s.staffRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.staffRepo.Delete(ctx, id); err != nil {
		return nil, err
	}
	return staff, nil
}
