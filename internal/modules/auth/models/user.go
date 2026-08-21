package models

import (
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	RoleID           uint           `json:"role_id,omitempty" gorm:"index"`
	Role             rolemodel.Role `json:"role"            gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name             string         `json:"name"            gorm:"type:varchar(100);index"`
	Email            string         `json:"email"           gorm:"type:varchar(150);not null"`
	Mobile           int            `json:"mobile"          gorm:"type:int;uniqueIndex"`
	Password         string         `json:"-"               gorm:"type:varchar(255);not null"`
	Status           bool           `json:"status"          gorm:"default:true;index"`
	UserType         uint8          `json:"user_type"       gorm:"index"`
	ProfilePicture   string         `json:"profile_picture" gorm:"type:varchar(500)"`
	CountryCode      string         `json:"country_code"    gorm:"type:char(5)"`
	Dob              string         `json:"dob"             gorm:"type:date"`
	DeviceToken      string         `json:"device_token"    gorm:"type:varchar(255)"`
	DeviceType       string         `json:"device_type"     gorm:"type:varchar(50)"`
	IsEmailVerified  bool           `json:"is_email_verified" gorm:"default:false"`
	IsMobileVerified bool           `json:"is_mobile_verified" gorm:"default:false"`
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
