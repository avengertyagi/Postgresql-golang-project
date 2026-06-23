package repositories

import (
	"github.com/akshit_tyagi/postgresql_project/internal/config"
	"github.com/akshit_tyagi/postgresql_project/internal/models"
)

func FindByEmail(email string) (*models.User, error) {
	var admin models.User
	err := config.DB.Where("email = ?", email).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func FindByID(id uint) (*models.User, error) {
	var admin models.User
	err := config.DB.First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func SaveToken(pat *models.PersonalAccessToken) error {
	return config.DB.Create(pat).Error
}

func FindTokenByHash(tokenHash string) (*models.PersonalAccessToken, error) {
	var pat models.PersonalAccessToken
	err := config.DB.
		Where("token_hash = ?", tokenHash).
		First(&pat).Error
	if err != nil {
		return nil, err
	}
	return &pat, nil
}
func RevokeRefreshToken(tokenHash string) error {
	return config.DB.
		Model(&models.PersonalAccessToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked", true).Error
}

func RevokeAllUserTokens(userID uint) error {
	return config.DB.
		Model(&models.PersonalAccessToken{}).
		Where("user_id = ? AND revoked = false", userID).
		Update("revoked", true).Error
}

func AssignRole(user *models.User, role *models.Role) error {
	return config.DB.Model(user).Association("Roles").Append(role)
}
