package repositories

import (
	"context"

	tempusermodel "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	"gorm.io/gorm"
)

type TempUserRepository interface {
	Create(ctx context.Context, user *tempusermodel.TempUser) error
	FindByEmail(ctx context.Context, email string) (*tempusermodel.TempUser, error)
	FindByMobile(ctx context.Context, mobile int) (*tempusermodel.TempUser, error)
	FindByID(ctx context.Context, id uint) (*tempusermodel.TempUser, error)
	Delete(ctx context.Context, user *tempusermodel.TempUser) error
}

type tempUserRepository struct {
	db *gorm.DB
}

func NewTempUserRepository(db *gorm.DB) TempUserRepository {
	return &tempUserRepository{db: db}
}

func (r *tempUserRepository) Create(ctx context.Context, user *tempusermodel.TempUser) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *tempUserRepository) FindByEmail(ctx context.Context, email string) (*tempusermodel.TempUser, error) {
	var user tempusermodel.TempUser
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *tempUserRepository) FindByMobile(ctx context.Context, mobile int) (*tempusermodel.TempUser, error) {
	var user tempusermodel.TempUser
	err := r.db.WithContext(ctx).Where("mobile = ?", mobile).First(&user).Error
	return &user, err
}

func (r *tempUserRepository) FindByID(ctx context.Context, id uint) (*tempusermodel.TempUser, error) {
	var user tempusermodel.TempUser
	err := r.db.WithContext(ctx).First(&user, id).Error
	return &user, err
}

func (r *tempUserRepository) Delete(ctx context.Context, user *tempusermodel.TempUser) error {
	return r.db.WithContext(ctx).Delete(user).Error
}
