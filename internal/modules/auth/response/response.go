package response

import (
	"github.com/akshit_tyagi/postgresql_project/internal/modules/auth/response/user"
)

type TokenResponse struct {
	AccessToken  string            `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string            `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType    string            `json:"token_type" example:"Bearer"`
	ExpiresIn    int               `json:"expires_in" example:"3600"`
	User         user.UserResponse `json:"user"`
}

type TokenRefreshResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType   string `json:"token_type" example:"Bearer"`
	ExpiresIn   int    `json:"expires_in" example:"3600"`
}
