package models

import (
	"gorm.io/gorm"
)

type TempUser struct {
	gorm.Model
	Name        string `json:"name"            gorm:"type:varchar(100);index"`
	Email       string `json:"email"           gorm:"type:varchar(150);not null"`
	Mobile      int    `json:"mobile"          gorm:"type:int;uniqueIndex"`
	Password    string `json:"-"              gorm:"type:varchar(255);not null"`
	UserType    uint8  `json:"user_type"       gorm:"index"`
	CountryCode string `json:"country_code"    gorm:"type:char(5)"`
	Dob         string `json:"dob"             gorm:"type:date"`
	DeviceToken string `json:"device_token"    gorm:"type:varchar(255)"`
	DeviceType  string `json:"device_type"     gorm:"type:varchar(50)"`
}
