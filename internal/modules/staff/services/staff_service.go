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

func List(ctx context.Context, q dto.PaginationQuery) ([]usermodel.User, int64, error) {
	return staffrepo.FindAll(ctx, q)
}

func Create(ctx context.Context, req dto.StaffRequest) (*usermodel.User, error) {
	name := strings.TrimSpace(req.Name)
	exists, err := staffrepo.ExistsByMobile(ctx, req.Mobile, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, staffconstants.StaffAlreadyExists
	}
	role, err := rolerepo.FindByID(ctx, req.RoleID)
	if err != nil || role == nil {
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
	if err := staffrepo.Create(ctx, staff); err != nil {
		return nil, err
	}
	return staff, nil
}

func GetByID(ctx context.Context, id uint) (*usermodel.User, error) {
	return staffrepo.FindByID(ctx, id)
}

func Update(ctx context.Context, id uint, req dto.StaffRequest) (*usermodel.User, error) {
	staff, err := staffrepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	exists, err := staffrepo.ExistsByMobile(ctx, req.Mobile, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, staffconstants.StaffAlreadyExists
	}
	staff.Name = name
	role, err := rolerepo.FindByID(ctx, req.RoleID)
	if err != nil || role == nil {
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
	if err := staffrepo.Update(ctx, staff); err != nil {
		return nil, err
	}

	return staff, nil
}

func UpdateStatus(ctx context.Context, id uint) (*usermodel.User, error) {
	staff, err := staffrepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	staff.Status = !staff.Status
	if err := staffrepo.Update(ctx, staff); err != nil {
		return nil, err
	}
	return staff, nil
}

func Delete(ctx context.Context, id uint) (*usermodel.User, error) {
	staff, err := staffrepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := staffrepo.Delete(ctx, id); err != nil {
		return nil, err
	}
	return staff, nil
}
