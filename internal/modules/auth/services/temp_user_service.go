package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	authconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/constants"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/dto"
	tempusermodel "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	authrepo "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/repositories"
	tempuserrepo "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(ctx context.Context, req dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	existingUser, _ := authrepo.FindByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, authconstants.EmailAlreadyExists
	}
	existingUser, _ = authrepo.FindByMobile(ctx, req.Mobile)
	if existingUser != nil {
		return nil, authconstants.MobileAlreadyExists
	}
	if existingTemp, tempErr := tempuserrepo.FindTempUserByEmail(ctx, req.Email); tempErr == nil && existingTemp != nil {
		return nil, authconstants.EmailAlreadyExists
	} else if tempErr != nil && !errors.Is(tempErr, gorm.ErrRecordNotFound) {
		return nil, tempErr
	}
	if existingTemp, tempErr := tempuserrepo.FindTempUserByMobile(ctx, req.Mobile); tempErr == nil && existingTemp != nil {
		return nil, authconstants.MobileAlreadyExists
	} else if tempErr != nil && !errors.Is(tempErr, gorm.ErrRecordNotFound) {
		return nil, tempErr
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &tempusermodel.TempUser{
		Name:        req.Name,
		Email:       req.Email,
		Mobile:      req.Mobile,
		CountryCode: req.CountryCode,
		Dob:         req.Dob,
		Password:    string(hashedPassword),
		UserType:    constants.UserRole,
	}
	if err := tempuserrepo.CreateTempUser(ctx, user); err != nil {
		slog.Error("error creating user", "error", err)
		return nil, err
	}
	response := &dto.UserRegisterResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Mobile:      user.Mobile,
		CountryCode: user.CountryCode,
		Dob:         user.Dob,
	}
	return response, nil
}
