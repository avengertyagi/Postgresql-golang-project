package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	authconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/constants"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/dto"
	"github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	authrepo "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/repositories"
)

type UserService interface {
	Register(ctx context.Context, req dto.UserRegisterRequest) (*dto.UserResponse, error)
}

type userService struct {
	repo authrepo.AuthRepository
}

func NewUserService(r authrepo.AuthRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) Register(ctx context.Context, req dto.UserRegisterRequest) (*dto.UserResponse, error) {
	existingUser, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, constants.UserNotFound) {
		slog.Error("error checking existing user by email", "error", err)
		return nil, err
	}
	if existingUser != nil {
		return nil, authconstants.EmailAlreadyExists
	}
	existingUser, err = s.repo.FindByMobile(ctx, req.Mobile)
	if err != nil && !errors.Is(err, constants.UserNotFound) {
		slog.Error("error checking existing user by mobile", "error", err)
		return nil, err
	}
	if existingUser != nil {
		return nil, authconstants.MobileAlreadyExists
	}
	user := &models.User{
		Name:        req.Name,
		Email:       req.Email,
		Mobile:      req.Mobile,
		CountryCode: req.CountryCode,
		Dob:         req.Dob,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		slog.Error("error creating user", "error", err)
		return nil, err
	}
	response := &dto.UserResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Mobile:      user.Mobile,
		CountryCode: user.CountryCode,
		Dob:         user.Dob,
	}
	return response, nil
}
