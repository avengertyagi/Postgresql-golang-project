package models

import (
	permissionmodel "github.com/akshit_tyagi/postgresql_project/internal/modules/permission/models"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string                       `json:"name"            gorm:"index;type:varchar(100);uniqueIndex;not null"`
	Permissions []permissionmodel.Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	Status      bool                         `json:"status"          gorm:"index;default:true"`
}
