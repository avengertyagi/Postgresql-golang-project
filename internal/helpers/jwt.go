package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	UserID      uint     `json:"-"`
	Email       string   `json:"-"`
	Role        string   `json:"-"`
	Guard       string   `json:"-"`
	Permissions []string `json:"-"`
	EncData     string   `json:"enc,omitempty"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	UserID  uint   `json:"-"`
	TokenID string `json:"-"`
	EncData string `json:"enc,omitempty"`
	jwt.RegisteredClaims
}

type encryptedAccessTokenPayload struct {
	UserID      uint     `json:"id"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Guard       string   `json:"guard"`
	Permissions []string `json:"permissions"`
}

type encryptedRefreshTokenPayload struct {
	UserID  uint   `json:"id"`
	TokenID string `json:"jti"`
}

func encryptTokenPayload(secret string, payload any) (string, error) {
	key := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plaintext, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

func decryptTokenPayload(secret, payload string, target any) error {
	key := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	ciphertext, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return err
	}
	if len(ciphertext) < gcm.NonceSize() {
		return fmt.Errorf("invalid encrypted token payload")
	}
	nonce, encryptedData := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return err
	}
	return json.Unmarshal(plaintext, target)
}

func GenerateAccessToken(userID uint, email, role, guard string, permissions []string) (string, error) {
	expiryMinutes, err := strconv.Atoi(os.Getenv("JWT_ACCESS_EXPIRY_MINUTES"))
	if err != nil || expiryMinutes == 0 {
		if os.Getenv("JWT_ACCESS_EXPIRY_MINUTES") != "" && err != nil {
			slog.Warn("jwt: invalid JWT_ACCESS_EXPIRY_MINUTES, defaulting to 60")
		}
		expiryMinutes = 60
	}
	secret := os.Getenv("JWT_ACCESS_SECRET")
	encryptedPayload, err := encryptTokenPayload(secret, encryptedAccessTokenPayload{
		UserID:      userID,
		Email:       email,
		Role:        role,
		Guard:       guard,
		Permissions: permissions,
	})
	if err != nil {
		return "", err
	}
	claims := AccessTokenClaims{
		EncData: encryptedPayload,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryMinutes) * time.Minute)),
			Issuer:    "auth-system",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID uint, tokenID string) (string, error) {
	expiryDays, err := strconv.Atoi(os.Getenv("JWT_REFRESH_EXPIRY_DAYS"))
	if err != nil || expiryDays == 0 {
		if os.Getenv("JWT_REFRESH_EXPIRY_DAYS") != "" && err != nil {
			slog.Warn("jwt: invalid JWT_REFRESH_EXPIRY_DAYS, defaulting to 30")
		}
		expiryDays = 30
	}
	secret := os.Getenv("JWT_REFRESH_SECRET")
	encryptedPayload, err := encryptTokenPayload(secret, encryptedRefreshTokenPayload{
		UserID:  userID,
		TokenID: tokenID,
	})
	if err != nil {
		return "", err
	}
	claims := RefreshTokenClaims{
		EncData: encryptedPayload,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryDays) * 24 * time.Hour)),
			Issuer:    "auth-system",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
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
	if claims.EncData != "" {
		var payload encryptedAccessTokenPayload
		if err := decryptTokenPayload(os.Getenv("JWT_ACCESS_SECRET"), claims.EncData, &payload); err != nil {
			return nil, fmt.Errorf("decrypt access token payload: %w", err)
		}
		claims.UserID = payload.UserID
		claims.Email = payload.Email
		claims.Role = payload.Role
		claims.Guard = payload.Guard
		claims.Permissions = payload.Permissions
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
	if claims.EncData != "" {
		var payload encryptedRefreshTokenPayload
		if err := decryptTokenPayload(os.Getenv("JWT_REFRESH_SECRET"), claims.EncData, &payload); err != nil {
			return nil, fmt.Errorf("decrypt refresh token payload: %w", err)
		}
		claims.UserID = payload.UserID
		claims.TokenID = payload.TokenID
	}
	return claims, nil
}

func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", h)
}
