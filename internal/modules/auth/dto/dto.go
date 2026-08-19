package dto

import "time"

type AdminLoginRequest struct {
	Email    string `json:"email" example:"admin@example.com"`
	Password string `json:"password" example:"password123"`
}
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIs..."`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIs..."`
}

type AdminResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	UserType     uint8  `json:"user_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type UserResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	UserType     uint8  `json:"user_type"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}
type TokenResponse struct {
	AccessToken  string       `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string       `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType    string       `json:"token_type" example:"Bearer"`
	ExpiresIn    int          `json:"expires_in" example:"3600"`
	User         UserResponse `json:"user"`
}

type TokenRefreshResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string `json:"token_type" example:"Bearer"`
	ExpiresIn   int    `json:"expires_in" example:"3600"`
}

type ProfileResponse struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	UserType       uint8     `json:"user_type"`
	Status         bool      `json:"status"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UpdateProfileRequest struct {
	Name           string `json:"name"            example:"John Doe"`
	ProfilePicture string `json:"profile_picture" example:"https://example.com/pic.jpg"`
}

type AdminUpdateProfileRequest struct {
	Name  string `form:"name"`
	Email string `form:"email"`
}
