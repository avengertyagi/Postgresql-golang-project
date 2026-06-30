package role

import (
	"time"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	permissionmodel "github.com/akshit_tyagi/postgresql_project/internal/models/permission"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID          uint                         `json:"id"              gorm:"primaryKey;autoIncrement"`
	TenantID    uint                         `json:"tenant_id"       gorm:"type:bigint;default:null"`
	Name        string                       `json:"name"            gorm:"type:varchar(100);uniqueIndex;not null"`
	Permissions []permissionmodel.Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	Status      bool                         `json:"status"          gorm:"default:true"`
	CreatedAt   time.Time                    `json:"created_at"`
	UpdatedAt   time.Time                    `json:"updated_at"`
	DeletedAt   gorm.DeletedAt               `gorm:"index" json:"deleted_at,omitempty"`
}

type RoleRequest struct {
	Name          string `json:"name"`
	PermissionIDs []uint `json:"permission"`
}

type RoleListRequest struct {
	CurrentPage int    `form:"currentPage"`
	PerPage     int    `form:"perPage"`
	Search      string `form:"search"`
	SortBy      string `form:"sortBy"`
	SortOrder   string `form:"sortOrder"`
}

type Pagination struct {
	CurrentPage int   `json:"currentPage"`
	PerPage     int   `json:"perPage"`
	Total       int64 `json:"total"`
	LastPage    int   `json:"lastPage"`
}

type RoleListAPIResponse struct {
	Status     bool       `json:"status"`
	StatusCode int        `json:"statusCode"`
	Message    string     `json:"message"`
	Data       []Role     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type RoleAPIResponse struct {
	Status     bool   `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       Role   `json:"data"`
}

type RoleListResponse struct {
	Data        []Role `json:"data"`
	CurrentPage int    `json:"currentPage"`
	PerPage     int    `json:"perPage"`
	Total       int64  `json:"total"`
	LastPage    int    `json:"lastPage"`
}

func (r *Role) SyncPermissions(permissionIDs []uint) error {
	var permissions []permissionmodel.Permission
	if err := config.DB.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}
	return config.DB.Model(r).Association("Permissions").Replace(permissions)
}
