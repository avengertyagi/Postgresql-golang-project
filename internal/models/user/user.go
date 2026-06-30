package user

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID             uint      `json:"id"              gorm:"primaryKey;autoIncrement"`
	TenantID       uint      `json:"tenant_id"       gorm:"type:bigint;default:null"`
	RoleID         uint      `json:"role_id"         gorm:"type:bigint;default:null"`
	Name           string    `json:"name"            gorm:"type:varchar(100);default:null"`
	Email          string    `json:"email"           gorm:"type:varchar(150);uniqueIndex;not null"`
	Password       string    `json:"-"               gorm:"type:varchar(255);not null"`
	Status         bool      `json:"status"          gorm:"default:true"`
	UserType       uint8     `json:"user_type" gorm:"default:1"`
	ProfilePicture string    `json:"profile_picture" gorm:"type:varchar(500);default:null"`
	DeviceToken    string    `json:"device_token"    gorm:"type:varchar(255);default:null"`
	DeviceType     string    `json:"device_type"     gorm:"type:varchar(50);default:null"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserLoginRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DeviceToken string `json:"device_token"`
	DeviceType  string `json:"device_type"`
}

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

type UserSignUpRequest struct {
	Name     string `json:"name" example:"John Doe"`
	Email    string `json:"email" example:"john.doe@example.com"`
	Password string `json:"password" example:"Password123"`
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

type UserResponse struct {
	ID    uint   `json:"id" example:"1"`
	Name  string `json:"name" example:"Admin User"`
	Email string `json:"email" example:"admin@example.com"`
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

type AdminAPIResponse struct {
	Status     bool          `json:"status"`
	StatusCode int           `json:"statusCode"`
	Message    string        `json:"message"`
	Data       AdminResponse `json:"data"`
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) HashPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
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
