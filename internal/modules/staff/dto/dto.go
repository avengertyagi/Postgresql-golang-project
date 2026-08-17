package dto

type PaginationQuery struct {
	Page    int    `form:"page,default=1" binding:"min=1"`
	Limit   int    `form:"limit,default=10" binding:"min=1,max=100"`
	Search  string `form:"search"`
	Status  string `form:"status" binding:"omitempty,oneof=draft published archived"`
	SortBy  string `form:"sort_by,default=created_at" binding:"omitempty,oneof=created_at title updated_at"`
	SortDir string `form:"sort_dir,default=desc" binding:"omitempty,oneof=asc desc"`
}

type StaffRequest struct {
	RoleID   uint   `json:"role_id" binding:"required"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Mobile   uint16 `json:"mobile"`
	Password string `json:"password"`
}
