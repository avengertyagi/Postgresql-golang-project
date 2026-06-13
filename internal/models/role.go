package models

import (
	"time"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
)

func init() {
	Register(&Role{})
}

type Role struct {
	ID        uint      `json:"id"              gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name"            gorm:"type:varchar(100);uniqueIndex;not null"`
	Status    bool      `json:"status"          gorm:"default:true"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RoleRequest struct {
	Name       string   `json:"name" gorm:"type:varchar(100);default:null"`
	Permission []string `json:"permission" gorm:"type:varchar(100);default:null"`
}

func SaveRole(role *Role) error {
	return config.DB.Create(role).Error
}

func FindRoleByName(name string) (*Role, error) {
	var role Role
	err := config.DB.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Role) SyncPermissions(permissions []Permission) error {
	return config.DB.Model(r).Association("Permissions").Replace(permissions)
}
