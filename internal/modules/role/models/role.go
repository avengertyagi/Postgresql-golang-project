package role

import (
	"github.com/akshit_tyagi/postgresql_project/internal/common"
	permissionmodel "github.com/akshit_tyagi/postgresql_project/internal/modules/permission/models"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string                       `json:"name"            gorm:"index;type:varchar(100);uniqueIndex;not null"`
	Permissions []permissionmodel.Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	Status      bool                         `json:"status"          gorm:"index;default:true"`
}

type RoleAPIResponse struct {
	ID          uint                         `json:"id"`
	Name        string                       `json:"name"`
	Permissions []permissionmodel.Permission `json:"permissions"`
	Status      bool                         `json:"status"`
}

type RoleRequest struct {
	Name          string `json:"name"`
	PermissionIDs []uint `json:"permission"`
}

type RoleListRequest struct {
	Page      int    `form:"page" json:"page"`
	Limit     int    `form:"limit" json:"limit"`
	Search    string `form:"search" json:"search"`
	SortBy    string `form:"sortBy" json:"sortBy"`
	SortOrder string `form:"sortOrder" json:"sortOrder"`
}

type RoleListAPIResponse struct {
	Data       []Role            `json:"data"`
	Pagination common.Pagination `json:"pagination"`
}
