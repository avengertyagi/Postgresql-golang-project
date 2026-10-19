package repositories

import (
	"context"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	usermodel "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	personalaccesstokenmodel "github.com/akshit_tyagi/postgresql_project/internal/modules/personalaccesstoken/models"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
)

func FindByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	var admin usermodel.User
	err := config.DB.WithContext(ctx).Model(&usermodel.User{}).Preload("Role.Permissions").Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func FindByMobile(ctx context.Context, mobile int) (*usermodel.User, error) {
	var admin usermodel.User
	err := config.DB.WithContext(ctx).Model(&usermodel.User{}).Preload("Role.Permissions").Where("mobile = ?", mobile).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func FindByID(ctx context.Context, id uint) (*usermodel.User, error) {
	var admin usermodel.User
	err := config.DB.WithContext(ctx).Model(&usermodel.User{}).Preload("Role.Permissions").Where("id = ?", id).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func Create(ctx context.Context, user *usermodel.User) error {
	return config.DB.WithContext(ctx).Create(user).Error
}

func Update(ctx context.Context, user *usermodel.User) error {
	return config.DB.WithContext(ctx).Save(user).Error
}

func SaveToken(ctx context.Context, pat *personalaccesstokenmodel.PersonalAccessToken) error {
	return config.DB.WithContext(ctx).Create(pat).Error
}

func FindTokenByHash(ctx context.Context, tokenHash string) (*personalaccesstokenmodel.PersonalAccessToken, error) {
	var pat personalaccesstokenmodel.PersonalAccessToken
	err := config.DB.WithContext(ctx).Model(&personalaccesstokenmodel.PersonalAccessToken{}).
		Where("token_hash = ?", tokenHash).
		First(&pat).Error
	if err != nil {
		return nil, err
	}
	return &pat, nil
}

func RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	return config.DB.WithContext(ctx).Model(&personalaccesstokenmodel.PersonalAccessToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked", true).Error
}

func RevokeAllUserTokens(ctx context.Context, userID uint) error {
	return config.DB.WithContext(ctx).Model(&personalaccesstokenmodel.PersonalAccessToken{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error
}

func AssignRole(ctx context.Context, user *usermodel.User, role *rolemodel.Role) error {
	user.RoleID = role.ID
	return config.DB.WithContext(ctx).Model(user).Update("role_id", role.ID).Error
}

func GetUserPermissions(ctx context.Context, user *usermodel.User) []string {
	var permissions []string
	for _, perm := range user.Role.Permissions {
		permissions = append(permissions, perm.Name)
	}
	return permissions
}
