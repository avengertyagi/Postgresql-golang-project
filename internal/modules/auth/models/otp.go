package models

import "time"

type OTP struct {
	ID         uint      `gorm:"primaryKey"`
	TempUserID uint      `gorm:"uniqueIndex;not null"`
	EmailCode  string    `gorm:"type:varchar(255);not null"`
	MobileCode string    `gorm:"type:varchar(255);not null"`
	ExpiresAt  time.Time `gorm:"not null"`
}
