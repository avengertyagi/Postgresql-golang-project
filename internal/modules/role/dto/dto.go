package dto

import (
	"github.com/akshit_tyagi/postgresql_project/internal/common"
	permissionmodel "github.com/akshit_tyagi/postgresql_project/internal/modules/permission/models"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
)

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

type PaginationQuery struct {
	Page    int    `form:"page,default=1" binding:"min=1"`
	Limit   int    `form:"limit,default=10" binding:"min=1,max=100"`
	Search  string `form:"search"`
	Status  string `form:"status" binding:"omitempty,oneof=draft published archived"`
	SortBy  string `form:"sort_by,default=created_at" binding:"omitempty,oneof=created_at title updated_at"`
	SortDir string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type RoleListAPIResponse struct {
	Data       []rolemodel.Role  `json:"data"`
	Pagination common.Pagination `json:"pagination"`
}
