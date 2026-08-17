package helpers

import (
	"crypto/sha256"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	UserID      uint     `json:"id"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Guard       string   `json:"guard"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID  uint   `json:"id"`
	TokenID string `json:"jti"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uint, email, role, guard string, permissions []string) (string, error) {
	expiryMinutes, err := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRY_MINUTES"))
	if err != nil || expiryMinutes == 0 {
		if os.Getenv("JWT_ACCESS_EXPIRY_MINUTES") != "" && err != nil {
			slog.Warn("jwt: invalid JWT_ACCESS_EXPIRY_MINUTES, defaulting to 60")
		}
		expiryMinutes = 60
	}
	claims := AccessTokenClaims{
		UserID:      userID,
		Email:       email,
		Role:        role,
		Guard:       guard,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryMinutes) * time.Minute)),
			Issuer:    "auth-system",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
}

func GenerateRefreshToken(userID uint, tokenID string) (string, error) {
	expiryDays, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRY_DAYS"))
	if err != nil || expiryDays == 0 {
		if os.Getenv("JWT_REFRESH_EXPIRY_DAYS") != "" && err != nil {
			slog.Warn("jwt: invalid JWT_REFRESH_EXPIRY_DAYS, defaulting to 30")
		}
		expiryDays = 30
	}
	claims := RefreshTokenClaims{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryDays) * 24 * time.Hour)),
			Issuer:    "auth-system",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
}

func ParseAccessToken(tokenStr string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid access token")
	}
	return claims, nil
}

func ParseRefreshToken(tokenStr string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}
	return claims, nil
}

func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}
