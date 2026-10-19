package repositories

import (
	"context"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	tempusermodel "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
)

func CreateTempUser(ctx context.Context, user *tempusermodel.TempUser) error {
	return config.DB.WithContext(ctx).Create(user).Error
}

func FindTempUserByEmail(ctx context.Context, email string) (*tempusermodel.TempUser, error) {
	var user tempusermodel.TempUser
	err := config.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func FindTempUserByMobile(ctx context.Context, mobile int) (*tempusermodel.TempUser, error) {
	var user tempusermodel.TempUser
	err := config.DB.WithContext(ctx).Where("mobile = ?", mobile).First(&user).Error
	return &user, err
}

func FindTempUserByID(ctx context.Context, id uint) (*tempusermodel.TempUser, error) {
	var user tempusermodel.TempUser
	err := config.DB.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func DeleteTempUser(ctx context.Context, user *tempusermodel.TempUser) error {
	return config.DB.WithContext(ctx).Delete(user).Error
}
