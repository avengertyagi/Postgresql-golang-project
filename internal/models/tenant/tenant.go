package tenant

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	gorm.Model
	ID             uint      `json:"id"              gorm:"primaryKey;autoIncrement"`
	Name           string    `json:"name"            gorm:"type:varchar(100);uniqueIndex;not null"`
	Email          string    `json:"email"           gorm:"type:varchar(100);not null"`
	Password       string    `json:"password"        gorm:"type:varchar(100);default:null"`
	Mobile         string    `json:"mobile"          gorm:"type:char(15); not null"`
	Address        string    `json:"address"         gorm:"type:mediumtext();default:null"`
	ProfilePicture string    `json:"profile_picture" gorm:"type:varchar(100);default:null"`
	PlanExpire     time.Time `json:"plan_expire"     gorm:"type:date"`
	Status         bool      `json:"status"          gorm:"default:true"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type TenantRequest struct {
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Password       string    `json:"password"`
	Mobile         string    `json:"mobile"`
	Address        string    `json:"address"`
	ProfilePicture string    `json:"profile_picture"`
	PlanExpire     time.Time `json:"plan_expire"`
}

type TenantListRequest struct {
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

type TenantListAPIResponse struct {
	Status     bool       `json:"status"`
	StatusCode int        `json:"statusCode"`
	Message    string     `json:"message"`
	Data       []Tenant   `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type TenantAPIResponse struct {
	Status     bool   `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       Tenant `json:"data"`
}

type TenantListResponse struct {
	Data        []Tenant `json:"data"`
	CurrentPage int      `json:"currentPage"`
	PerPage     int      `json:"perPage"`
	Total       int64    `json:"total"`
	LastPage    int      `json:"lastPage"`
}
