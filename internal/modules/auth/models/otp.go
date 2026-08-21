package models

import "time"

type OTP struct {
	ID         uint      `gorm:"primaryKey"`
	TempUserID uint      `gorm:"uniqueIndex;not null"`
	EmailCode  int       `gorm:"type:varchar(255);not null"`
	MobileCode int       `gorm:"type:varchar(255);not null"`
	ExpiresAt  time.Time `gorm:"not null"`
}
