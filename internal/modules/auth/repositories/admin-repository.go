package repositories

import (
	"context"

	usermodel "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	personalaccesstokenmodel "github.com/akshit_tyagi/postgresql_project/internal/modules/personalaccesstoken/models"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	FindByEmail(ctx context.Context, email string) (*usermodel.User, error)
	FindByID(ctx context.Context, id uint) (*usermodel.User, error)
	SaveToken(ctx context.Context, pat *personalaccesstokenmodel.PersonalAccessToken) error
	FindTokenByHash(ctx context.Context, tokenHash string) (*personalaccesstokenmodel.PersonalAccessToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllUserTokens(ctx context.Context, userID uint) error
	AssignRole(ctx context.Context, user *usermodel.User, role *rolemodel.Role) error
	GetUserPermissions(ctx context.Context, user *usermodel.User) []string
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) FindByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	var admin usermodel.User
	err := r.db.WithContext(ctx).Model(&usermodel.User{}).Preload("Role.Permissions").Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *authRepository) FindByID(ctx context.Context, id uint) (*usermodel.User, error) {
	var admin usermodel.User
	err := r.db.Model(&usermodel.User{}).Preload("Role.Permissions").Where("id = ?", id).First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (r *authRepository) SaveToken(ctx context.Context, pat *personalaccesstokenmodel.PersonalAccessToken) error {
	return r.db.WithContext(ctx).Create(pat).Error
}

func (r *authRepository) FindTokenByHash(ctx context.Context, tokenHash string) (*personalaccesstokenmodel.PersonalAccessToken, error) {
	var pat personalaccesstokenmodel.PersonalAccessToken
	err := r.db.WithContext(ctx).Model(&personalaccesstokenmodel.PersonalAccessToken{}).
		Where("token_hash = ?", tokenHash).
		First(&pat).Error
	if err != nil {
		return nil, err
	}
	return &pat, nil
}
func (r *authRepository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	return r.db.WithContext(ctx).Model(&personalaccesstokenmodel.PersonalAccessToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked", true).Error
}

func (r *authRepository) RevokeAllUserTokens(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Model(&personalaccesstokenmodel.PersonalAccessToken{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error
}

func (r *authRepository) AssignRole(ctx context.Context, user *usermodel.User, role *rolemodel.Role) error {
	user.RoleID = role.ID
	return r.db.WithContext(ctx).Model(user).Update("role_id", role.ID).Error
}

func (r *authRepository) GetUserPermissions(ctx context.Context, user *usermodel.User) []string {
	var permissions []string
	for _, perm := range user.Role.Permissions {
		permissions = append(permissions, perm.Name)
	}
	return permissions
}
