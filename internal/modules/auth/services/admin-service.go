package admin

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	helpers "github.com/akshit_tyagi/postgresql_project/internal/helpers"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/dto"
	usermodel "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/models"
	authrepo "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/repositories"
	personalaccesstokenmodel "github.com/akshit_tyagi/postgresql_project/internal/modules/personalaccesstoken/models"
	"github.com/google/uuid"
)

type AdminService interface {
	Login(ctx context.Context, req dto.AdminLoginRequest) (*dto.AdminResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshToken(ctx context.Context, refreshToken string) (*dto.TokenRefreshResponse, error)
	GetProfile(ctx context.Context, userID uint) (*usermodel.ProfileResponse, error)
}
type authService struct {
	repo authrepo.AuthRepository
}

func NewAuthService(repo authrepo.AuthRepository) AdminService {
	return &authService{repo: repo}
}

func (s *authService) Login(ctx context.Context, req dto.AdminLoginRequest) (*dto.AdminResponse, error) {
	admin, err := s.repo.FindByEmail(ctx, req.Email)
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
	permissions := s.repo.GetUserPermissions(ctx, admin)
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
	if err := s.repo.SaveToken(ctx, pat); err != nil {
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

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	if _, err := helpers.ParseRefreshToken(refreshToken); err != nil {
		return constants.SessionNotFound
	}
	tokenHash := helpers.HashToken(refreshToken)
	pat, err := s.repo.FindTokenByHash(ctx, tokenHash)
	if err != nil {
		return constants.SessionNotFound
	}
	if pat.Revoked {
		return constants.SessionAlreadyRevoked
	}
	if time.Now().After(pat.ExpiresAt) {
		return constants.SessionExpired
	}
	if err := s.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return constants.SomethingWentWrong
	}
	return nil
}

func (s *authService) RefreshToken(ctx context.Context, rawRefreshToken string) (*dto.TokenRefreshResponse, error) {
	claims, err := helpers.ParseRefreshToken(rawRefreshToken)
	if err != nil {
		return nil, constants.SessionNotFound
	}
	pat, err := s.repo.FindTokenByHash(ctx, helpers.HashToken(rawRefreshToken))
	if err != nil {
		return nil, errors.New("Session not found")
	}
	if pat.Revoked {
		return nil, errors.New("Session has been revoked")
	}
	if time.Now().After(pat.ExpiresAt) {
		return nil, errors.New("Session has expired")
	}
	admin, err := s.repo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, constants.UserNotFound
	}
	permissions := s.repo.GetUserPermissions(ctx, admin)
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

func (s *authService) GetProfile(ctx context.Context, userID uint) (*usermodel.ProfileResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, constants.NotFound
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
