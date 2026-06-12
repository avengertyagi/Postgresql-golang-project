package service

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/akshit_tyagi/postgresql_project/src/constants"
	authmodel "github.com/akshit_tyagi/postgresql_project/src/models"
	"github.com/akshit_tyagi/postgresql_project/src/utils"
	"github.com/google/uuid"
)

func Login(req authmodel.AdminLoginRequest) (*authmodel.AdminResponse, error) {
	admin, err := authmodel.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New(constants.INVALID_CREDENTIALS)
	}
	isSuperAdmin := admin.UserType == constants.SUPER_ADMIN_ROLE
	if !isSuperAdmin && !admin.Status {
		return nil, errors.New(constants.INACTIVE_ACCOUNT)
	}
	if !admin.CheckPassword(req.Password) {
		return nil, errors.New(constants.INVALID_CREDENTIALS)
	}
	userID := admin.ID
	accessToken, err := utils.GenerateAccessToken(
		userID,
		admin.Email,
		strconv.Itoa(int(admin.UserType)),
		constants.ADMIN_GUARD,
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
	pat := &authmodel.PersonalAccessToken{
		UserID:    admin.ID,
		TokenHash: utils.HashToken(rawRefreshToken),
		Name:      "admin-session",
		Revoked:   false,
		ExpiresAt: time.Now().Add(time.Duration(expiryDays) * 24 * time.Hour),
	}
	if err := authmodel.SaveToken(pat); err != nil {
		return nil, errors.New("Failed to create session")
	}
	accessExpiryMinutes, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRY_MINUTES"))
	if accessExpiryMinutes == 0 {
		accessExpiryMinutes = 60
	}
	return &authmodel.AdminResponse{
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
		return errors.New(constants.SESSION_NOT_FOUND)
	}
	tokenHash := utils.HashToken(refreshToken)
	pat, err := authmodel.FindTokenByHash(tokenHash)
	if err != nil {
		return errors.New(constants.SESSION_NOT_FOUND)
	}
	if pat.Revoked {
		return errors.New(constants.SESSION_ALREADY_REVOKED)
	}
	if time.Now().After(pat.ExpiresAt) {
		return errors.New(constants.SESSION_EXPIRED)
	}
	if err := authmodel.RevokeRefreshToken(tokenHash); err != nil {
		return errors.New("Failed to revoke session")
	}
	return nil
}

func RefreshToken(rawRefreshToken string) (*authmodel.TokenRefreshResponse, error) {
	claims, err := utils.ParseRefreshToken(rawRefreshToken)
	if err != nil {
		return nil, errors.New(constants.SESSION_NOT_FOUND)
	}
	pat, err := authmodel.FindTokenByHash(utils.HashToken(rawRefreshToken))
	if err != nil {
		return nil, errors.New("Session not found")
	}
	if pat.Revoked {
		return nil, errors.New("Session has been revoked")
	}
	if time.Now().After(pat.ExpiresAt) {
		return nil, errors.New("Session has expired")
	}
	admin, err := authmodel.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New(constants.USER_NOT_FOUND)
	}
	newAccessToken, err := utils.GenerateAccessToken(
		claims.UserID,
		admin.Email,
		strconv.Itoa(int(admin.UserType)),
		constants.ADMIN_GUARD,
	)
	if err != nil {
		return nil, errors.New("Failed to generate access token")
	}
	accessExpiryMinutes, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRY_MINUTES"))
	if accessExpiryMinutes == 0 {
		accessExpiryMinutes = 60
	}
	return &authmodel.TokenRefreshResponse{
		AccessToken: newAccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   accessExpiryMinutes * 60,
	}, nil
}

func GetProfile(userID uint) (*authmodel.ProfileResponse, error) {
	user, err := authmodel.FindByID(userID)
	if err != nil {
		return nil, errors.New(constants.NOT_FOUND)
	}
	return &authmodel.ProfileResponse{
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
