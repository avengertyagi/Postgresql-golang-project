package dto

import (
	"github.com/akshit_tyagi/postgresql_project/internal/common"
	thememodel "github.com/akshit_tyagi/postgresql_project/internal/modules/theme/models"
)

type ThemeRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type PaginationQuery struct {
	Page    int    `form:"page,default=1" binding:"min=1"`
	Limit   int    `form:"limit,default=10" binding:"min=1,max=100"`
	Search  string `form:"search"`
	Status  string `form:"status" binding:"omitempty,oneof=true false"`
	SortBy  string `form:"sort_by,default=created_at" binding:"omitempty,oneof=name created_at updated_at"`
	SortDir string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type ThemeListAPIResponse struct {
	Data       []thememodel.Theme `json:"data"`
	Pagination common.Pagination  `json:"pagination"`
}
