package repositories

import (
	"github.com/akshit_tyagi/postgresql_project/internal/config"
	usermodel "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	personalaccesstokenmodel "github.com/akshit_tyagi/postgresql_project/internal/modules/personalaccesstoken/models"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
)

func FindByEmail(email string) (*usermodel.User, error) {
	var admin usermodel.User
	err := config.DB.Model(&usermodel.User{}).Preload("Role.Permissions").Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func FindByID(id uint) (*usermodel.User, error) {
	var admin usermodel.User
	err := config.DB.Model(&usermodel.User{}).Preload("Role.Permissions").Where("id = ?", id).First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func SaveToken(pat *personalaccesstokenmodel.PersonalAccessToken) error {
	return config.DB.Create(pat).Error
}

func FindTokenByHash(tokenHash string) (*personalaccesstokenmodel.PersonalAccessToken, error) {
	var pat personalaccesstokenmodel.PersonalAccessToken
	err := config.DB.Model(&personalaccesstokenmodel.PersonalAccessToken{}).
		Where("token_hash = ?", tokenHash).
		First(&pat).Error
	if err != nil {
		return nil, err
	}
	return &pat, nil
}
func RevokeRefreshToken(tokenHash string) error {
	return config.DB.Model(&personalaccesstokenmodel.PersonalAccessToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked", true).Error
}

func RevokeAllUserTokens(userID uint) error {
	return config.DB.Model(&personalaccesstokenmodel.PersonalAccessToken{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error
}

func AssignRole(user *usermodel.User, role *rolemodel.Role) error {
	user.RoleID = role.ID
	return config.DB.Model(user).Update("role_id", role.ID).Error
}

func GetUserPermissions(user *usermodel.User) []string {
	var permissions []string
	for _, perm := range user.Role.Permissions {
		permissions = append(permissions, perm.Name)
	}
	return permissions
}
