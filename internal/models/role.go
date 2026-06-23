package models

import (
	"time"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
)

func init() {
	Register(&Role{})
}

type Role struct {
	ID          uint         `json:"id"              gorm:"primaryKey;autoIncrement"`
	Name        string       `json:"name"            gorm:"type:varchar(100);uniqueIndex;not null"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	Status      bool         `json:"status"          gorm:"default:true"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type RoleRequest struct {
	Name          string `json:"name" gorm:"type:varchar(100);default:null"`
	PermissionIDs []uint `json:"permission"`
}

func (r *Role) SyncPermissions(permissionIDs []uint) error {
	var permissions []Permission
	if err := config.DB.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}
	return config.DB.Model(r).Association("Permissions").Replace(permissions)
}
