package services

import (
	"context"
	"errors"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	helpers "github.com/akshit_tyagi/postgresql_project/internal/helpers"
	authconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/constants"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/dto"
	authrepo "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/repositories"
	personalaccesstokenmodel "github.com/akshit_tyagi/postgresql_project/internal/modules/personalaccesstoken/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const dummyBcryptHash = "$2a$10$CwTycUXWue0Thq9StjUM0uJ8Kge0G9wORUiC5MzHzWMbc8dV/xhHy"

func Login(ctx context.Context, req dto.AdminLoginRequest) (*dto.AdminResponse, error) {
	admin, err := authrepo.FindByEmail(ctx, req.Email)
	if err != nil {
		_ = bcrypt.CompareHashAndPassword([]byte(dummyBcryptHash), []byte(req.Password))
		return nil, authconstants.InvalidCredentials
	}
	isSuperAdmin := admin.UserType == constants.SuperAdminRole
	if !isSuperAdmin && !admin.Status {
		return nil, authconstants.InactiveAccount
	}
	if !admin.CheckPassword(req.Password) {
		return nil, authconstants.InvalidCredentials
	}
	userID := admin.ID
	permissions := authrepo.GetUserPermissions(ctx, admin)
	accessToken, err := helpers.GenerateAccessToken(
		userID,
		admin.Email,
		strconv.Itoa(int(admin.UserType)),
		constants.AdminGuard,
		permissions,
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
	if err := authrepo.SaveToken(ctx, pat); err != nil {
		return nil, errors.New("Failed to create session")
	}
	accessExpiryMinutes, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRY_MINUTES"))
	if accessExpiryMinutes == 0 {
		accessExpiryMinutes = 60
	}
	return &dto.AdminResponse{
		ID:           userID,
		Name:         admin.Name,
		Email:        admin.Email,
		UserType:     admin.UserType,
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		ExpiresIn:    accessExpiryMinutes * 60,
	}, nil
}

func Logout(ctx context.Context, refreshToken string) error {
	if _, err := helpers.ParseRefreshToken(refreshToken); err != nil {
		return authconstants.SessionNotFound
	}
	tokenHash := helpers.HashToken(refreshToken)
	pat, err := authrepo.FindTokenByHash(ctx, tokenHash)
	if err != nil {
		return authconstants.SessionNotFound
	}
	if pat.Revoked {
		return authconstants.SessionAlreadyRevoked
	}
	if time.Now().After(pat.ExpiresAt) {
		return authconstants.SessionExpired
	}
	if err := authrepo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return constants.SomethingWentWrong
	}
	return nil
}

func RefreshToken(ctx context.Context, rawRefreshToken string) (*dto.TokenRefreshResponse, error) {
	claims, err := helpers.ParseRefreshToken(rawRefreshToken)
	if err != nil {
		return nil, authconstants.SessionNotFound
	}
	pat, err := authrepo.FindTokenByHash(ctx, helpers.HashToken(rawRefreshToken))
	if err != nil {
		return nil, errors.New("Session not found")
	}
	if pat.Revoked {
		return nil, errors.New("Session has been revoked")
	}
	if time.Now().After(pat.ExpiresAt) {
		return nil, errors.New("Session has expired")
	}
	admin, err := authrepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, constants.UserNotFound
	}
	permissions := authrepo.GetUserPermissions(ctx, admin)
	newAccessToken, err := helpers.GenerateAccessToken(
		claims.UserID,
		admin.Email,
		strconv.Itoa(int(admin.UserType)),
		constants.AdminGuard,
		permissions,
	)
	if err != nil {
		return nil, errors.New("Failed to generate access token")
	}
	accessExpiryMinutes, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRY_MINUTES"))
	if accessExpiryMinutes == 0 {
		accessExpiryMinutes = 60
	}
	return &dto.TokenRefreshResponse{
		AccessToken: newAccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   accessExpiryMinutes * 60,
	}, nil
}

func GetProfile(ctx context.Context, userID uint) (*dto.ProfileResponse, error) {
	user, err := authrepo.FindByID(ctx, userID)
	if err != nil {
		return nil, constants.NotFound
	}
	return &dto.ProfileResponse{
		ID:             user.ID,
		Name:           helpers.StringOrNA(&user.Name),
		Email:          helpers.StringOrNA(&user.Email),
		UserType:       user.UserType,
		Status:         user.Status,
		ProfilePicture: helpers.StringOrNA(&user.ProfilePicture),
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}, nil
}

func UpdateProfile(ctx context.Context, userID uint, req dto.AdminUpdateProfileRequest, profilePictureFile interface{}) (*dto.ProfileResponse, error) {
	user, err := authrepo.FindByID(ctx, userID)
	if err != nil {
		return nil, constants.NotFound
	}
	if name := strings.TrimSpace(req.Name); name != "" {
		user.Name = name
	}
	if email := strings.TrimSpace(req.Email); email != "" {
		user.Email = email
	}
	if profilePictureFile != nil {
		file, ok := profilePictureFile.(*multipart.FileHeader)
		if !ok {
			return nil, errors.New("invalid file format")
		}
		result, err := helpers.UploadSingleImage(ctx, file, "profile-pictures")
		if err != nil {
			return nil, err
		}
		user.ProfilePicture = result.URL
	}

	if err := authrepo.Update(ctx, user); err != nil {
		return nil, constants.SomethingWentWrong
	}
	return GetProfile(ctx, userID)
}
