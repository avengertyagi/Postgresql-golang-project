package admin

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	helpers "github.com/akshit_tyagi/postgresql_project/internal/helpers"
	personalaccesstokenmodel "github.com/akshit_tyagi/postgresql_project/internal/models/personalaccesstoken"
	usermodel "github.com/akshit_tyagi/postgresql_project/internal/models/user"
	repositories "github.com/akshit_tyagi/postgresql_project/internal/repositories/user"
	userrepo "github.com/akshit_tyagi/postgresql_project/internal/repositories/user"
	"github.com/google/uuid"
)

func Login(req usermodel.AdminLoginRequest) (*usermodel.AdminResponse, error) {
	admin, err := userrepo.FindByEmail(req.Email)
	if err != nil {
		return nil, constants.InvalidCredentials
	}
	isSuperAdmin := admin.UserType == constants.SuperAdminRole
	if !isSuperAdmin && !admin.Status {
		return nil, constants.InactiveAccount
	}
	if !admin.CheckPassword(req.Password) {
		return nil, constants.InvalidCredentials
	}
	userID := admin.ID
	accessToken, err := helpers.GenerateAccessToken(
		userID,
		admin.Email,
		strconv.Itoa(int(admin.UserType)),
		constants.AdminGuard,
	)
	if err != nil {
		return nil, errors.New("Failed to generate access token")
	}
	jti := uuid.New().String()
	rawRefreshToken, err := helpers.GenerateRefreshToken(userID, jti)
	if err != nil {
		return nil, errors.New("Failed to generate refresh token")
	}
	expiryDays, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRY_DAYS"))
	if expiryDays == 0 {
		expiryDays = 30
	}
	pat := &personalaccesstokenmodel.PersonalAccessToken{
		UserID:    admin.ID,
		TokenHash: helpers.HashToken(rawRefreshToken),
		Name:      "admin-session",
		Revoked:   false,
		ExpiresAt: time.Now().Add(time.Duration(expiryDays) * 24 * time.Hour),
	}
	if err := userrepo.SaveToken(pat); err != nil {
		return nil, errors.New("Failed to create session")
	}
	accessExpiryMinutes, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRY_MINUTES"))
	if accessExpiryMinutes == 0 {
		accessExpiryMinutes = 60
	}
	return &usermodel.AdminResponse{
		ID:           userID,
		Name:         admin.Name,
		Email:        admin.Email,
		UserType:     admin.UserType,
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		ExpiresIn:    accessExpiryMinutes * 60,
	}, nil
}

func Logout(refreshToken string) error {
	if _, err := helpers.ParseRefreshToken(refreshToken); err != nil {
		return errors.New(constants.SessionNotFound)
	}
	tokenHash := helpers.HashToken(refreshToken)
	pat, err := repositories.FindTokenByHash(tokenHash)
	if err != nil {
		return errors.New(constants.SessionNotFound)
	}
	if pat.Revoked {
		return errors.New(constants.SessionAlreadyRevoked)
	}
	if time.Now().After(pat.ExpiresAt) {
		return errors.New(constants.SessionExpired)
	}
	if err := repositories.RevokeRefreshToken(tokenHash); err != nil {
		return errors.New("Failed to revoke session")
	}
	return nil
}

func RefreshToken(rawRefreshToken string) (*usermodel.TokenRefreshResponse, error) {
	claims, err := helpers.ParseRefreshToken(rawRefreshToken)
	if err != nil {
		return nil, errors.New(constants.SessionNotFound)
	}
	pat, err := repositories.FindTokenByHash(helpers.HashToken(rawRefreshToken))
	if err != nil {
		return nil, errors.New("Session not found")
	}
	if pat.Revoked {
		return nil, errors.New("Session has been revoked")
	}
	if time.Now().After(pat.ExpiresAt) {
		return nil, errors.New("Session has expired")
	}
	admin, err := repositories.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New(constants.UserNotFound)
	}
	newAccessToken, err := helpers.GenerateAccessToken(
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
	return &usermodel.TokenRefreshResponse{
		AccessToken: newAccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   accessExpiryMinutes * 60,
	}, nil
}

func GetProfile(userID uint) (*usermodel.ProfileResponse, error) {
	user, err := repositories.FindByID(userID)
	if err != nil {
		return nil, errors.New(constants.NotFound)
	}
	return &usermodel.ProfileResponse{
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
