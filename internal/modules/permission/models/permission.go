package models

import (
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	Name      string `json:"name"       gorm:"type:varchar(100);uniqueIndex;not null"`
	GuardName string `json:"guard_name" gorm:"type:varchar(100);default:null"`
	Status    bool   `json:"status"     gorm:"default:true"`
}
