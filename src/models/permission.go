package models

import (
	"time"

	"github.com/akshit_tyagi/postgresql_project/src/config"
)

func init() {
	Register(&Permission{})
}

type Permission struct {
	ID        uint      `json:"id"              gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name"            gorm:"type:varchar(100);uniqueIndex;not null"`
	GuardName string    `json:"guard_name"      gorm:"type:varchar(100);default:null"`
	Status    bool      `json:"status"          gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func SavePermission(permission *Permission) error {
	return config.DB.Create(permission).Error
}

func FindPermissionByName(name string) (*Permission, error) {
	var permission Permission
	err := config.DB.Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}
