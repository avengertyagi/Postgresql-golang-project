package user

import (
	"time"

	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	RoleID         uint           `json:"role_id"         gorm:"type:bigint;default:null"`
	Role           rolemodel.Role `json:"role"            gorm:"foreignKey:role_id"`
	Name           string         `json:"name"            gorm:"index;type:varchar(100);default:null"`
	Email          string         `json:"email"           gorm:"index;type:varchar(150);uniqueIndex;not null"`
	Password       string         `json:"-"               gorm:"type:varchar(255);not null"`
	Status         bool           `json:"status"          gorm:"index;default:true"`
	UserType       uint8          `json:"user_type" gorm:"index;default:null"`
	ProfilePicture string         `json:"profile_picture" gorm:"type:varchar(500);default:null"`
	DeviceToken    string         `json:"device_token"    gorm:"type:varchar(255);default:null"`
	DeviceType     string         `json:"device_type"     gorm:"type:varchar(50);default:null"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIs..."`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIs..."`
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
