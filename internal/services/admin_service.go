package services

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	"github.com/akshit_tyagi/postgresql_project/internal/models"
	"github.com/akshit_tyagi/postgresql_project/internal/utils"
	"github.com/google/uuid"
)

func Login(req models.AdminLoginRequest) (*models.AdminResponse, error) {
	admin, err := models.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New(constants.InvalidCredentials)
	}
	isSuperAdmin := admin.UserType == constants.SuperAdminRole
	if !isSuperAdmin && !admin.Status {
		return nil, errors.New(constants.InactiveAccount)
	}
	if !admin.CheckPassword(req.Password) {
		return nil, errors.New(constants.InvalidCredentials)
	}
	userID := admin.ID
	accessToken, err := utils.GenerateAccessToken(
		userID,
		admin.Email,
		strconv.Itoa(int(admin.UserType)),
		constants.AdminGuard,
	)
	if err != nil {
		return nil, errors.New("Failed to generate access token")
	}
	jti := uuid.New().String()
	rawRefreshToken, err := utils.GenerateRefreshToken(userID, jti)
	if err != nil {
		return nil, errors.New("Failed to generate refresh token")
	}
	expiryDays, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRY_DAYS"))
	if expiryDays == 0 {
		expiryDays = 30
	}
	pat := &models.PersonalAccessToken{
		UserID:    admin.ID,
		TokenHash: utils.HashToken(rawRefreshToken),
		Name:      "admin-session",
		Revoked:   false,
		ExpiresAt: time.Now().Add(time.Duration(expiryDays) * 24 * time.Hour),
	}
	if err := models.SaveToken(pat); err != nil {
		return nil, errors.New("Failed to create session")
	}
	accessExpiryMinutes, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRY_MINUTES"))
	if accessExpiryMinutes == 0 {
		accessExpiryMinutes = 60
	}
	return &models.AdminResponse{
		ID:           userID,
		Name:         admin.Name,
		Email:        admin.Email,
		UserType:     admin.UserType,
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    accessExpiryMinutes * 60,
	}, nil
}

func Logout(refreshToken string) error {
	if _, err := utils.ParseRefreshToken(refreshToken); err != nil {
		return errors.New(constants.SessionNotFound)
	}
	tokenHash := utils.HashToken(refreshToken)
	pat, err := models.FindTokenByHash(tokenHash)
	if err != nil {
		return errors.New(constants.SessionNotFound)
	}
	if pat.Revoked {
		return errors.New(constants.SessionAlreadyRevoked)
	}
	if time.Now().After(pat.ExpiresAt) {
		return errors.New(constants.SessionExpired)
	}
	if err := models.RevokeRefreshToken(tokenHash); err != nil {
		return errors.New("Failed to revoke session")
	}
	return nil
}

func RefreshToken(rawRefreshToken string) (*models.TokenRefreshResponse, error) {
	claims, err := utils.ParseRefreshToken(rawRefreshToken)
	if err != nil {
		return nil, errors.New(constants.SessionNotFound)
	}
	pat, err := models.FindTokenByHash(utils.HashToken(rawRefreshToken))
	if err != nil {
		return nil, errors.New("Session not found")
	}
	if pat.Revoked {
		return nil, errors.New("Session has been revoked")
	}
	if time.Now().After(pat.ExpiresAt) {
		return nil, errors.New("Session has expired")
	}
	admin, err := models.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New(constants.UserNotFound)
	}
	newAccessToken, err := utils.GenerateAccessToken(
		claims.UserID,
		admin.Email,
		strconv.Itoa(int(admin.UserType)),
		constants.AdminGuard,
	)
	if err != nil {
		return nil, errors.New("Failed to generate access token")
	}
	accessExpiryMinutes, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRY_MINUTES"))
	if accessExpiryMinutes == 0 {
		accessExpiryMinutes = 60
	}
	return &models.TokenRefreshResponse{
		AccessToken: newAccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   accessExpiryMinutes * 60,
	}, nil
}

func GetProfile(userID uint) (*models.ProfileResponse, error) {
	user, err := models.FindByID(userID)
	if err != nil {
		return nil, errors.New(constants.NotFound)
	}
	return &models.ProfileResponse{
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		UserType:       user.UserType,
		Status:         user.Status,
		ProfilePicture: user.ProfilePicture,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}, nil
}
